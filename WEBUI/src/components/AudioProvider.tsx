import * as React from 'react';
import { connect } from 'react-redux';
import { Dispatch, Action } from 'redux';
import { ButtonBase, Typography } from '@material-ui/core';
import PlayArrowRoundedIcon from '@material-ui/icons/PlayArrowRounded';
import SkipNextRoundedIcon from '@material-ui/icons/SkipNextRounded';
import PauseRoundedIcon from '@material-ui/icons/PauseRounded';
import styled, { keyframes } from 'styled-components';
import { Program } from '../Types/Main';
import { RootState } from '../Types';
import { mainActionCreator } from '../Actions/Main';

interface StateToProps {
    nowprog?: Program,
    queue: Array<Program>,
}
interface DispatchToProps {
    skip: () => void,
}

interface State {
    isPlay: boolean,
    isOpen: boolean,
    player: HTMLAudioElement,
    nowProg?: Program,
}
type Props = StateToProps & DispatchToProps;

class AudioProvider extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = {
            isPlay: false,
            isOpen: false,
            player: new Audio('https://object-storage.tyo1.conoha.io/v1/nc_92a6769609d54403bc799a178c136a31/radio/test'),
        };
        this.changePlay = this.changePlay.bind(this);
    }

    componentDidMount() {
        const { player } = this.state;
        const { skip } = this.props;
        player.addEventListener('ended', () => {
            const { queue } = this.props;
            if (queue.length > 0) {
                skip();
            }
        });
    }

    componentDidUpdate(prevProps: Props, prevState: State) {
        const { nowprog } = this.props;
        if (prevProps.nowprog?.id !== nowprog?.id) {
            const { player, isPlay, ...other } = this.state;
            player.src = nowprog.uri;
            // eslint-disable-next-line react/no-did-update-set-state
            this.setState({
                ...other,
                player,
                isPlay: false,
                nowProg: nowprog,
            });
            setTimeout(() => this.changePlay(), 5);
        }
    }

    changePlay() {
        const { isPlay, player, ...other } = this.state;
        if (other.nowProg === undefined) {
            return;
        }
        if (isPlay) {
            player.pause();
        } else {
            player.play();
        }
        this.setState({
            ...other,
            isPlay: !isPlay,
            player,
        });
    }

    render() {
        const { nowProg, isPlay } = this.state;
        const { skip, queue } = this.props;
        const isSkip = queue.length > 0;
        return (
            <Root>
                <TickerWrap>
                    <TickerMain>
                        <Title variant="h6">{(nowProg) ? nowProg.title : 'にじさんじpresentsだいたいにじさんじのらじお'}</Title>
                    </TickerMain>
                </TickerWrap>
                <Comment variant="body1">{(nowProg === undefined) ? '' : nowProg.recTimestamp}</Comment>
                <LongBox>
                    <CenterButtonBase onClick={this.changePlay}>
                        {(isPlay) ? <PauseRoundedIcon /> : <PlayArrowRoundedIcon /> }
                    </CenterButtonBase>
                </LongBox>
                <LongBox>
                    <CenterButtonBase disabled={!isSkip} onClick={() => skip()}>
                        {(isSkip) ? (
                            <SkipNextRoundedIcon />
                        ) : <SkipNextRoundedIcon color="disabled" />}
                    </CenterButtonBase>
                </LongBox>
            </Root>
        );
    }
}

const mapStateToProps = (state: RootState): StateToProps => {
    const { Main } = state;
    return {
        nowprog: Main.nowProgram,
        queue: Main.audioQueue,
    };
};

const mapDispatchToProps = (dispatch: Dispatch<Action>): DispatchToProps => ({
    skip: () => {
        dispatch(mainActionCreator.skip());
    },
});

export default connect(
    mapStateToProps,
    mapDispatchToProps,
)(AudioProvider);


const Root = styled.div`
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 12%;
    background-color: #f0f0f0;
    padding: 5px 10px 5px 5px;
    margin: 0px;
    display: grid;
    grid-auto-flow: column;
    grid-template-columns: 6fr 1fr 1fr;
    grid-template-rows: 1fr 1fr;
`;

const Title = styled(Typography)`
    margin: auto 5px;
    // text-overflow: ellipsis;
    // width: 95%;
    // overflow: hidden;
    // white-space: nowrap;
`;

const LongBox = styled.div`
    grid-row: 1 / 3;
    display: flex;
    justify-content: center;
    align-items: center;
`;

const CenterButtonBase = styled(ButtonBase)`
    width: 35px;
    height: 35px;
`;

const Comment = styled(Typography)`
    margin: 0px 5px;
`;

/* ticker */
const Ticker = keyframes`
    0% {
        transform: translate(0, 0);
        visibility: visible;
    }
    100% {
        transform: translate(-100%, 0);
    }
`;

const TickerMain = styled.div`
    display: inline-block;
    height: 2rem;
    line-height: 2rem;
    white-space: nowrap;
    box-sizing: content-box;
    animation: ${Ticker} 30s linear infinite;
    animation-delay: 3s;
`;
const TickerWrap = styled.div`
    overflow: hidden;
    height: 100%;
    box-sizing: content-box;
`;
