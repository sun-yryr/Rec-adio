import React from 'react';
import { Card, CardContent, CardActionArea } from '@material-ui/core';
import { Program } from '../Types/Main';
import AudioProvider from './AudioProvider';

interface Props extends Program {
    play: (url: string) => void
}

export const ResultCell = (props: Props) => {
    const { play, title, uri } = props;
    const clickhandler = () => {
        play(uri);
    };

    return (
        <Card>
            <CardActionArea onClick={clickhandler}>
                <CardContent>
                    <h2>{title}</h2>
                    <AudioProvider src={uri} />
                </CardContent>
            </CardActionArea>
        </Card>
    );
};
