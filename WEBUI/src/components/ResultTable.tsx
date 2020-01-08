import React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import { connect } from 'react-redux';
import { Dispatch, Action } from 'redux';
import styled from 'styled-components';
import { RootState } from '../Types';
import { Program } from '../Types/Main';
import { ResultCell } from './ResultCell';
import { mainActionCreator } from '../Actions/Main';

interface StateToProps {
    onFetch: boolean,
    data: Array<Program>,
}
interface DispatchToProps {
    addFront: (prog: Program) => void,
    addQueue: (prog: Program) => void,
}
type IProps = StateToProps & DispatchToProps & RouteComponentProps;

const ResultTable = (props: IProps) => {
    const {
        onFetch,
        data,
        addFront,
        addQueue,
    } = props;

    if (onFetch) { return <p>loading</p>; }
    return (
        <Grid>
            {data.map((prog) => (
                <GridItem>
                    <ResultCell key={prog.id} addFront={addFront} addQueue={addQueue} program={prog} />
                </GridItem>
            ))}
        </Grid>
    );
};

const mapStateToProps = (state: RootState): StateToProps => {
    const { Api } = state;
    return {
        onFetch: Api.onFetch,
        data: Api.data,
    };
};

const mapDispatchToProps = (dispatch: Dispatch<Action>): DispatchToProps => ({
    addFront: (prog: Program) => {
        dispatch(mainActionCreator.addFront(prog));
    },
    addQueue: (prog: Program) => {
        dispatch(mainActionCreator.addQueue(prog));
    },
});

export default connect(
    mapStateToProps,
    mapDispatchToProps,
)(withRouter(ResultTable));

const Grid = styled.div`
    display: grid;
    grid-template-rows: repeat(7, 1fr);
    height: 100%;
`;

const GridItem = styled.div`
    border-top: solid;
    border-bottom: solid;
    border-width: thin;
`;
