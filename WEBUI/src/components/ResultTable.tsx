import React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import { connect } from 'react-redux';
import { RootState } from '../Types';
import { Program } from '../Types/Main';
import { ResultCell } from './ResultCell';

interface StateToProps {
    onFetch: boolean,
    data: Array<Program>,
}
type IProps = StateToProps & RouteComponentProps;

const ResultTable = (props: IProps) => {
    const { onFetch, data } = props;
    const play = (url: string) => {
        const tmp = new Audio();
        tmp.play();
    };

    if (onFetch) { return <p>loading</p>; }
    return (
        <div>
            {data.map((prog) => <ResultCell key={prog.id} {...prog} play={play} />)}
        </div>
    );
};

const mapStateToProps = (state: RootState): StateToProps => {
    const { Api } = state;
    return {
        onFetch: Api.onFetch,
        data: Api.data,
    };
};

export default connect(
    mapStateToProps,
)(withRouter(ResultTable));
