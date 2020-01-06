import * as React from 'react';
import { connect } from 'react-redux';
import { Dispatch, Action } from 'redux';
import { Button } from '@material-ui/core';
import styled from 'styled-components';
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
        const { nowprog, queue } = this.props;
        if (prevProps.nowprog?.id !== nowprog?.id) {
            const { player, isPlay, ...other } = this.state;
            player.src = nowprog.uri;
            this.setState({
                ...other,
                player,
                isPlay: false,
                nowProg: nowprog,
            });
            if (true) {
                setTimeout(() => this.changePlay(), 5);
            }
        }
    }

    changePlay() {
        const { isPlay, player, ...other } = this.state;
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
        console.log(nowProg);
        return (
            <Root>
                <Title>{(nowProg === undefined) ? '' : nowProg.title}</Title>
                <Title>{(nowProg === undefined) ? '' : nowProg.recTimestamp}</Title>
                <LongBox>
                    <Button variant="outlined" onClick={this.changePlay}>{(isPlay) ? 'Stop' : 'Start'}</Button>
                </LongBox>
                <LongBox>
                    {(queue.length > 0) ? (
                        <Button variant="outlined" onClick={() => skip()}>skip</Button>
                    ) : <Button variant="outlined" disabled onClick={() => skip()}>skip</Button>}
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
    padding: 5px;
    margin: 0px;
    display: grid;
    grid-auto-flow: column;
    grid-template-columns: 3fr 1fr 1fr;
    grid-template-rows: 1fr 1fr;
`;

const Title = styled.p`
    margin: auto 5px;
    text-overflow: clip;
    width: 95%;
    overflow: auto;
    white-space: nowrap;
`;

const LongBox = styled.div`
    grid-row: 1 / 3;
    display: flex;
`;
