var express = require('express');
var router = express.Router();

var dotenv = require('dotenv');
dotenv.config();

const mariadb = require('mariadb');
const pool = mariadb.createPool({
     host: 'localhost',
     user: process.env.DBUSER,
     password: process.env.DBPASS,
     database: process.env.DBNAME,
     connectionLimit: 5
});

/*
    channelの検索 なくす
    channel一覧を取得                list/channels
    channelIdを指定して動画情報の取得   videos/channel/:Id
    -----------------
    グループ一覧を取得                 list/groups
    channel検索のオートコンプリート
    検索クエリによる動画情報の取得       videos?q
    グループの動画情報の取得            videos/group/:Id
*/

/*
    api/search ワード検索
        パラメータ q : 検索クエリ
                 matched : full | default
                          完全一致
                 isInfo : Boolean
*/

function oreore_mysql_escape(str) {
    return str.replace(/(["'`;])/g, `\\$1`);
}

/* チャンネル一覧の取得 */
router.get('/api/search', async function (req, res, next) {
    // qがない場合，エラー出して終了
    let q = req.query.q || '';
    if (q === '') {
        res.status = 400;
        res.json({Error: "params[q] is undefined."});
        return;
    }
    // dbにアクセス，Channelsから情報を取得
    let conn;
    let rows;
    try {
        conn = await pool.getConnection();
        rows = await conn.query("SELECT * FROM Programs WHERE `title` LIKE ? OR `pfm` LIKE ?", [q, q]);
        res.status = 200
        res.json(rows);
    } catch (err) {
        res.status = 400;
        res.json({Error: err});
        throw err;
    } finally {
        if (conn) conn.release(); //release to pool
    }
});

module.exports = router;
