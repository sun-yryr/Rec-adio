import React from '../../node_modules/@types/react';
import styled from '../../node_modules/@types/styled-components/ts3.7';
import { withRouter, RouteComponentProps } from '../../node_modules/@types/react-router-dom';
import { Button } from '../../node_modules/@material-ui/core';

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
