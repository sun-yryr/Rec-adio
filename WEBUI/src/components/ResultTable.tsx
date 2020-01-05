import React from 'react';
import styled from 'styled-components';
import { RouteComponentProps, withRouter } from 'react-router';
import { connect } from 'react-redux';
import { Program, RootState } from '../Types';

interface StateToProps {
    onFetch: boolean,
    data: Array<Program>,
}
type IProps = StateToProps & RouteComponentProps;

const ResultTable = (props: IProps) => {
    const { onFetch, data } = props;
    if (onFetch) { return <p>loading</p>; }
    return (
        <div>
            {data.map((prog) => <p>{prog.title}</p>)}
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
