import React from 'react';
import styled from 'styled-components';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import { connect } from 'react-redux';
import { Button, TextField } from '@material-ui/core';
import { RootState } from '../Types';
import { Dispatch, Action } from 'redux';
import { mainActionCreator } from '../Actions/Main';

interface Props extends RouteComponentProps {
    title: string
}
interface State {
    text: string,
}
interface StateToProps {
    pass: string
}
interface DispatchToProps {
    setPass: (pass: string) => void;
}
type IProps = Props & StateToProps & DispatchToProps;

class Top extends React.Component<IProps, State> {
    constructor(props: any) {
        super(props);
        this.state = { text: '' };
        this.change = this.change.bind(this);
        this.save = this.save.bind(this);
    }

    change(e: any) {
        const text: string = e.target.value;
        this.setState({ text });
    }
    save() {
        const { text } = this.state;
        if (text === '') {
            return;
        }
        const { setPass } = this.props;
        setPass(text);
    }

    render() {
        const { pass } = this.props;
        console.log(pass);
        return (
            <Root>
                {(pass === '') ? (
                    <div>
                        <TextField onChange={this.change} />
                        <button onClick={this.save}>保存</button>
                    </div>
                ) : null}
                <Button
                    variant="contained"
                    onClick={() => this.props.history.push('/search')}
                >
                    検索
                </Button>
            </Root>
        );
    }
};

const mapStateToProps = (state: RootState): StateToProps => {
    const { Main } = state;
    return {
        pass: Main.password,
    };
}

const mapDispatchToProps = (dispatch: Dispatch<Action>): DispatchToProps => ({
    setPass: (pass: string) => {
        dispatch(mainActionCreator.setPassword(pass));
    },
})

export default connect(
    mapStateToProps,
    mapDispatchToProps,
)(withRouter(Top));

const Root = styled.div`
    text-align: center;
`;
