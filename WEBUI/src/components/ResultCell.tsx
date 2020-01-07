import React from 'react';
import styled from 'styled-components';
import {
    Button,
    Typography,
} from '@material-ui/core';
import { Program } from '../Types/Main';

interface Props {
    addFront: (prog: Program) => void,
    addQueue: (prog: Program) => void,
    program: Program,
}

export const ResultCell = (props: Props) => {
    const { addFront, program, addQueue } = props;

    return (
        <CardContentGrid>
            <Title variant="body1">{program.title}</Title>
            <Text variant="caption">{program.recTimestamp}</Text>
            <ButtonSet>
                <Button variant="outlined" onClick={() => addFront(program)}>Play</Button>
                <Button variant="outlined" onClick={() => addQueue(program)}>Add</Button>
            </ButtonSet>
        </CardContentGrid>
    );
};

const CardContentGrid = styled.div`
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    grid-template-rows: 1fr 1fr;
    margin: 10px;
`;

const Title = styled(Typography)`
    grid-column: 1 / 4;
`;

const Text = styled(Typography)`
    grid-column: 1 / 3;
    margin-left: 10px;
`;

const ButtonSet = styled.div`
    grid-column: 3/ 4;
    display: grid;
    grid-template-columns: 1fr 1fr;
`;
