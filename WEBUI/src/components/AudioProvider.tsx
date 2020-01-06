import * as React from 'react';
import { Program } from '../Types/Main';

interface Props {
    src: string
}
interface State {
    isPlay: boolean,
    isOpen: boolean,
    player?: HTMLAudioElement,
    nowProg?: Program,
}

class AudioProvider extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = { 
            isPlay: false,
            isOpen: false,
        };
        this.play = this.play.bind(this);
    }

    play() {
        const { isPlay } = this.state;
        if (isPlay) {
            return;
        }
        this.setState({ isPlay: true });
        const { src } = this.props;
        const audio = new Audio('https://object-storage.tyo1.conoha.io/v1/nc_92a6769609d54403bc799a178c136a31/radio/test');
        audio.play();
    }

    render() {
        return (
            <button onClick={this.play}>start</button>
        );
    }
}

export default AudioProvider;
