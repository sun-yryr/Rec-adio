package recorder

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sun-yryr/recoto/internal/domain"
	"github.com/sun-yryr/recoto/internal/event/recording"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
)

const (
	radikoTimeLayout = "20060102150405"
	radikoAuthKey    = "bcd151073c03b352e1ef2fd66c32209da9ca0afa"

	// HTTP status codes.
	httpStatusOK = 200
)

// RadikoRecorder はRadikoをソースとして録音を行うrecorder。
type RadikoRecorder struct {
	logger    *zap.Logger
	radikoURL string
	areaID    string
}

// RadikoConfig はRadikoRecorderの設定。
type RadikoConfig struct {
	// RadikoURL はRadikoのプログラム情報を取得するURL。
	// デフォルトは "http://radiko.jp/v3/program/today/JP13.xml"
	RadikoURL string
	// AreaID はRadikoのエリアID。
	// デフォルトは "JP13" (東京)
	AreaID string
}

// NewRadikoRecorder はRadikoRecorderを生成するコンストラクタ。
func NewRadikoRecorder(logger *zap.Logger, config RadikoConfig) *RadikoRecorder {
	radikoURL := "http://radiko.jp/v3/program/today/JP13.xml"
	if config.RadikoURL != "" {
		radikoURL = config.RadikoURL
	}

	areaID := "JP13"
	if config.AreaID != "" {
		areaID = config.AreaID
	}

	return &RadikoRecorder{
		logger:    logger,
		radikoURL: radikoURL,
		areaID:    areaID,
	}
}

// GetSupportSource はRadikoRecorderがサポートするソースを返す。
func (r *RadikoRecorder) GetSupportSource() []domain.SourceKind {
	return []domain.SourceKind{domain.SourceKindRadiko}
}

// CheckAvailable はRadikoRecorderの有効性をチェックする。
func (r *RadikoRecorder) CheckAvailable() error {
	// ffmpegの存在を確認
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return errors.Wrap(err, "ffmpeg is not installed")
	}

	return nil
}

// GetName はRadikoRecorderの名前を返す。
func (r *RadikoRecorder) GetName() string {
	return "RadikoRecorder"
}

// Rec はRadikoをソースとして録音を行う。
func (r *RadikoRecorder) Rec(ctx context.Context, event *recording.RequestedEvent) error {
	stationID := event.Source.ID

	// Radikoの認証を行う
	authToken, err := r.authorization(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to authorize radiko")
	}

	// ストリームURLを取得
	streamURL := fmt.Sprintf("http://f-radiko.smartstream.ne.jp/%s/_definst_/simul-stream.stream/playlist.m3u8", stationID)
	m3u8URL, err := r.getM3U8URL(ctx, streamURL, authToken)
	if err != nil {
		return errors.Wrap(err, "failed to get m3u8 url")
	}

	var stdErrBuf bytes.Buffer

	// ffmpegコマンドを構築
	cmd := ffmpeg_go.
		Input(
			m3u8URL,
			ffmpeg_go.KwArgs{
				"headers": fmt.Sprintf("X-Radiko-Authtoken: %s\r\n", authToken),
				"t":       strconv.FormatInt(int64(math.Ceil(event.Duration.Seconds())), 10),
			},
		).
		Output(event.Output, ffmpeg_go.KwArgs{"acodec": "copy"}).
		WithErrorOutput(&stdErrBuf)

	// キャンセル付きのcontextを設定する
	cmd.Context = ctx
	// 処理の猶予時間を追加
	const gracePeriod = 10 * time.Second
	// int64を超える（290年を超える）録音時間に対する保護
	maxDuration := time.Duration(math.MaxInt64) - gracePeriod
	if event.Duration > maxDuration {
		return errors.New("recording duration is too long")
	}

	cmd = cmd.WithTimeout(event.Duration + gracePeriod)

	r.logger.Debug(
		"start recording",
		zap.String("recordingId", event.RecordingID),
		zap.String("stationId", stationID),
		zap.String("cmd", cmd.String()),
	)

	if err := cmd.Run(); err != nil {
		r.logger.Debug(
			"failed to run ffmpeg",
			zap.String("stderr", stdErrBuf.String()),
			zap.Error(err),
		)

		return errors.Wrap(err, "failed to run ffmpeg")
	}

	r.logger.Debug(
		"finish recording",
		zap.String("recordingId", event.RecordingID),
		zap.String("stationId", stationID),
	)

	return nil
}

// RadikoProgram はRadikoの番組情報を表す構造体。
type RadikoProgram struct {
	XMLName  xml.Name `xml:"radiko"`
	Stations []struct {
		ID    string `xml:"id,attr"`
		Name  string `xml:"name"`
		Progs struct {
			Progs []struct {
				FT    string `xml:"ft,attr"`
				TO    string `xml:"to,attr"`
				FTL   string `xml:"ftl,attr"`
				TOL   string `xml:"tol,attr"`
				Dur   string `xml:"dur,attr"`
				Title string `xml:"title"`
				Info  string `xml:"info"`
				Desc  string `xml:"desc"`
				Pfm   string `xml:"pfm"`
			} `xml:"prog"`
		} `xml:"progs"`
	} `xml:"stations>station"`
}

// SearchProgram はキーワードに一致する番組を検索する。
func (r *RadikoRecorder) SearchProgram(keywords []string) ([]map[string]interface{}, error) {
	if len(keywords) == 0 {
		return nil, nil
	}

	// キーワードを正規表現に変換
	var pattern string

	for i, keyword := range keywords {
		if i > 0 {
			pattern += "|"
		}

		pattern += regexp.QuoteMeta(keyword)
	}

	regexPattern, err := regexp.Compile("(" + pattern + ")")
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile regex")
	}

	// 番組情報を取得
	program, err := r.fetchProgram()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch program")
	}

	// 検索結果
	var results []map[string]interface{}

	// 番組を検索
	for _, station := range program.Stations {
		for _, prog := range station.Progs.Progs {
			match := false

			// タイトルを検索
			if regexPattern.MatchString(prog.Title) {
				match = true
			}

			// 情報を検索
			if !match && prog.Info != "" && regexPattern.MatchString(prog.Info) {
				match = true
			}

			// 出演者を検索
			if !match && prog.Pfm != "" && regexPattern.MatchString(prog.Pfm) {
				match = true
			}

			// 詳細を検索
			if !match && prog.Desc != "" && regexPattern.MatchString(prog.Desc) {
				match = true
			}

			if match {
				// 開始時刻をパース
				startTimeValue, err := time.Parse(radikoTimeLayout, prog.FT)
				if err != nil {
					r.logger.Warn(
						"failed to parse start time",
						zap.Error(err),
						zap.String("ft", prog.FT),
					)

					continue
				}

				// 番組時間をパース
				durSec, err := strconv.Atoi(prog.Dur)
				if err != nil {
					r.logger.Warn(
						"failed to parse duration",
						zap.Error(err),
						zap.String("dur", prog.Dur),
					)

					continue
				}

				// 結果に追加
				results = append(results, map[string]interface{}{
					"station": station.ID,
					"title":   strings.ReplaceAll(prog.Title, " ", "_"),
					"ft":      prog.FT,
					"DT_ft":   startTimeValue,
					"to":      prog.TO,
					"ftl":     prog.FTL,
					"tol":     prog.TOL,
					"dur":     durSec,
					"pfm":     strings.ReplaceAll(prog.Pfm, "，", ","),
					"info":    prog.Info,
				})
			}
		}
	}

	return results, nil
}

// authorization はRadikoの認証を行い、認証トークンを返す。
func (r *RadikoRecorder) authorization(ctx context.Context) (string, error) {
	auth1URL := "https://radiko.jp/v2/api/auth1"
	auth2URL := "https://radiko.jp/v2/api/auth2"

	// auth1リクエスト
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, auth1URL, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create auth1 request")
	}

	req.Header.Set("X-Radiko-App", "pc_html5")
	req.Header.Set("X-Radiko-App-Version", "0.0.1")
	req.Header.Set("X-Radiko-User", "sunyryr")
	req.Header.Set("X-Radiko-Device", "pc")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {

		return "", errors.Wrap(err, "failed to send auth1 request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.Warn("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != httpStatusOK {
		return "", errors.Wrap(
			errors.New("auth1 failed"),
			fmt.Sprintf("status code: %d", resp.StatusCode),
		)
	}

	authToken := resp.Header.Get("X-Radiko-Authtoken")
	keyLength := resp.Header.Get("X-Radiko-Keylength")
	keyOffset := resp.Header.Get("X-Radiko-Keyoffset")

	// キー長とオフセットを数値に変換
	keyLengthInt, err := strconv.Atoi(keyLength)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse key length")
	}

	keyOffsetInt, err := strconv.Atoi(keyOffset)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse key offset")
	}

	// 部分キーを生成
	partialKey := radikoAuthKey[keyOffsetInt : keyOffsetInt+keyLengthInt]
	authKey := base64.StdEncoding.EncodeToString([]byte(partialKey))

	// auth2リクエスト
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, auth2URL, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create auth2 request")
	}

	req.Header.Set("X-Radiko-Authtoken", authToken)
	req.Header.Set("X-Radiko-Partialkey", authKey)
	req.Header.Set("X-Radiko-User", "sunyryr")
	req.Header.Set("X-Radiko-Device", "pc")

	resp, err = client.Do(req)

	if err != nil {
		return "", errors.Wrap(err, "failed to send auth2 request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.Warn("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != httpStatusOK {
		return "", errors.Wrap(
			errors.New("auth2 failed"),
			fmt.Sprintf("status code: %d", resp.StatusCode),
		)
	}

	return authToken, nil
}

// getM3U8URL はストリームURLからM3U8 URLを取得する。
func (r *RadikoRecorder) getM3U8URL(ctx context.Context, streamURL, authToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, streamURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create m3u8 request")
	}

	req.Header.Set("X-Radiko-Authtoken", authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to send m3u8 request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.Warn("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != httpStatusOK {
		return "", errors.Wrap(
			errors.New("m3u8 request failed"),
			fmt.Sprintf("status code: %d", resp.StatusCode),
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read m3u8 response")
	}

	// M3U8 URLを正規表現で抽出
	re := regexp.MustCompile(`^https?://.+m3u8$`)
	matches := re.FindAllString(string(body), -1)

	if len(matches) == 0 {
		return "", errors.New("no m3u8 url found in response")
	}

	return matches[0], nil
}

// fetchProgram はRadikoの番組情報を取得する。
func (r *RadikoRecorder) fetchProgram() (*RadikoProgram, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.radikoURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create program request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get program")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.Warn("failed to close response body", zap.Error(err))
		}
	}()

	var program RadikoProgram
	if err := xml.NewDecoder(resp.Body).Decode(&program); err != nil {
		return nil, errors.Wrap(err, "failed to decode program")
	}

	return &program, nil
}
