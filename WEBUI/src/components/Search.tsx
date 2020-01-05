import React from 'react';
import styled from 'styled-components';
import { RouteComponentProps, withRouter } from 'react-router';
import { TextField, Button } from '@material-ui/core';
import { connect } from 'react-redux';
import { ThunkDispatch } from 'redux-thunk';
import { RootState, RootActions } from '../Types';
import { apiActionCreator } from '../Actions/Api';

interface State {
    keyword: string
}

interface DispatchToProps {
    search: (query: string) => void,
}
type IProps = DispatchToProps & RouteComponentProps;

class Search extends React.Component<IProps, State> {
    constructor(props: IProps) {
        super(props);
        this.state = { keyword: '' };
        this.input = this.input.bind(this);
        this.search = this.search.bind(this);
    }

    input(e: any) {
        const keyword: string = e.target.value;
        this.setState({ keyword });
    }

    search() {
        const { search, history } = this.props;
        const { keyword } = this.state;
        search(keyword);
        history.push('/result');
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
                <Button variant="contained" onClick={this.search}>検索</Button>
            </Root>
        );
    }
}

const mapDispatchToProps = (dispatch: ThunkDispatch<RootState, undefined, RootActions>): DispatchToProps => ({
    search: (query: string) => {
        dispatch(apiActionCreator.getData(query));
    },
});

export default connect(
    mapDispatchToProps,
)(withRouter(Search));

const Root = styled.div`
    text-align: center;
`;
