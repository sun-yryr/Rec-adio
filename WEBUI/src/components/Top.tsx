import React from 'react';
import styled from 'styled-components';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import { Button } from '@material-ui/core';

interface Props extends RouteComponentProps {
    title: string
}

const Top = (props: Props) => {
    const { title } = props;
    return (
        <Root>
            <p>{title}</p>
            <Button
                variant="contained"
                onClick={() => props.history.push('/search')}
            >
                検索
            </Button>
        </Root>
    );
};

export default withRouter(Top);

const Root = styled.div`
    text-align: center;
`;
