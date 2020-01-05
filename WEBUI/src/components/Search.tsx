import React, { ChangeEvent } from '../../node_modules/@types/react';
import styled from '../../node_modules/@types/styled-components/ts3.7';
import { RouteComponentProps, withRouter } from '../../node_modules/@types/react-router';
import { TextField, Button } from '../../node_modules/@material-ui/core';

interface State {
    keyword: string
}

class Search extends React.Component<RouteComponentProps, State> {
    constructor(props: RouteComponentProps) {
        super(props);
        this.state = { keyword: '' };
        this.input = this.input.bind(this);
    }

    input(e: any) {
        const keyword: string = e.target.value;
        this.setState({ keyword });
    }

    render() {
        const { keyword } = this.state;
        return (
            <Root>
                <TextField
                    id="outlined-full-width"
                    style={{
                        margin: 8,
                        width: '80%',
                    }}
                    placeholder="search keyword"
                    margin="normal"
                    InputLabelProps={{
                        shrink: true,
                    }}
                    variant="outlined"
                    value={keyword}
                    onChange={this.input}
                />
                <Button variant="contained">検索</Button>
            </Root>
        );
    }
}

export default withRouter(Search);

const Root = styled.div`
    text-align: center;
`;
