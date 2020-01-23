import React from 'react';
import styled from 'styled-components';
import PlayArrowRoundedIcon from '@material-ui/icons/PlayArrowRounded';
import QueuePlayNextRoundedIcon from '@material-ui/icons/QueuePlayNextRounded';
import {
    Button,
    Typography,
    ButtonBase,
} from '@material-ui/core';
import moment from 'moment';
import { Program } from '../Types/Main';

interface Props {
    addFront: (prog: Program) => void,
    addQueue: (prog: Program) => void,
    program: Program,
}

export const ResultCell = (props: Props) => {
    const { addFront, program, addQueue } = props;
    const wrapAddFront = () => {
        addFront(program);
    };
    const wrapAddQueue = () => {
        addQueue(program);
    };
    const formatDate = moment(program.recTimestamp).format('YYYY年MM月DD日');

    return (
        <CardContentGrid>
            <Title variant="body1">{program.title}</Title>
            <Text variant="caption">{formatDate}</Text>
            <ButtonSet>
                <ExtendBBase
                    onClick={wrapAddFront}
                    style={{
                        padding: '5px 10px',
                    }}
                >
                    <PlayArrowRoundedIcon />
                </ExtendBBase>
                <ExtendBBase
                    onClick={wrapAddQueue}
                    style={{
                        padding: '5px 10px',
                    }}
                >
                    <QueuePlayNextRoundedIcon />
                </ExtendBBase>
            </ButtonSet>
        </CardContentGrid>
    );
};

const CardContentGrid = styled.div`
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    grid-template-rows: 1fr 1fr;
    margin: 5px 5px 5px 10px;
`;

const Title = styled(Typography)`
    grid-column: 1 / 4;
`;

const Text = styled(Typography)`
    grid-column: 1 / 3;
    margin-left: 10px;
`;

const ButtonSet = styled.div`
    grid-column: 3 / 4;
    display: grid;
    justify-content: end;
    grid-auto-flow: column;
`;

const ExtendBBase = styled(ButtonBase)`
    margin: 5px 10px;
`;
