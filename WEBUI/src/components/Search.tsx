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

interface StateToProps {}
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
                <GridItem>
                    <TextField
                        id="outlined-full-width"
                        style={{
                            padding: 8,
                            margin: 0,
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
                    <SearchButton variant="contained" onClick={this.search}>検索</SearchButton>
                </GridItem>
            </Root>
        );
    }
}

const mapStateToProps = (): StateToProps => ({});

const mapDispatchToProps = (dispatch: ThunkDispatch<RootState, undefined, RootActions>): DispatchToProps => ({
    search: (query: string) => {
        dispatch(apiActionCreator.getData(query));
    },
});

export default connect(
    mapStateToProps,
    mapDispatchToProps,
)(withRouter(Search));

const Root = styled.div`
    display: grid;
    grid-auto-flow: column;
    grid-template-rows: 80px auto;
    height: 100%;
`;

const GridItem = styled.div`
    display: grid;
    grid-template-columns: 3fr 1fr;
    align-content: center;
`;

const SearchButton = styled(Button)`
    height: 55px;
    padding: 0;
    width: 90%;
    margin: auto 0px;
`;
