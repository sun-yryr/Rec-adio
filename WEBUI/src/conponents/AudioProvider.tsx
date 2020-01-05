import * as React from '../../node_modules/@types/react';
import { Url } from 'url';

interface Props {
    src: string
}
interface State {
    isPlay: boolean
}

class AudioProvider extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = { isPlay: false };
        this.play = this.play.bind(this);
    }

    play() {
        const { isPlay } = this.state;
        if (isPlay) {
            return;
        }
        this.setState({ isPlay: true });
        const { src } = this.props;
        const audio = new Audio(src);
        audio.play();
    }

    render() {
        return (
            <button onClick={this.play}>start</button>
        );
    }
}

export default AudioProvider;
