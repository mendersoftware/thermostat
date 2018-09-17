import React, { Component } from 'react';
import io from 'socket.io-client';
import nipple from 'nipplejs';
import TimerMixin from 'react-timer-mixin';


const socket = io.connect(location.protocol + '//' + document.domain + ':' + location.port);

class Timer extends React.Component {
    constructor(props) {
        super(props);
    }
    componentDidMount() {
        this.timer = TimerMixin.setInterval(() => {
            this.props.onUpdate();
        }, this.props.interval);
    }
    componentWillUnmount() {
        clearTimeout(this.timer);
    }
    render() {
        return null;
    }
}

class Joystick extends React.Component {
    constructor(props) {
        super(props);
        this.joystick = null;
    }

    componentDidMount() {
        let self = this;
        let options =  {
            zone: this.joystickField,
            mode: 'static',
            position: {
                left: '10%',
                top: '20%'
            },
            color: 'red',
            threshold: '0.2',
            fadeTime: '500',
        }
        this.joystick = nipple.create(options);
        this.joystick.on('move', function(evt, data) {
            let angle = (data.direction == null) ? null : data.direction.angle
            self.props.onMove(data.force, data.angle.degree, angle)
            
        }).on('end', function(evt, data) {
            self.props.onStop();
        });
    }
    
    render() {
        return (
            <div> 
                <div ref={c => this.joystickField = c} />
                <ul>
                    <li>force: {this.props.force} </li>
                    <li>angle: {this.props.angle} </li>
                </ul>
            </div>
        );
    }
}

class App extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            joystick: {
                force: 0,
                angle: 0,
                dir: null
            },
            interval: 500
        }
    }

    sendStateToSocket() {
        console.log('will emit data ',
            JSON.stringify(JSON.stringify(
                {
                    force: this.state.joystick.force,
                    direction: this.state.joystick.dir, 
                    interval: this.state.interval
                })
            )
        );

        socket.emit('move', 
            JSON.stringify(
                {
                    force: this.state.joystick.force,
                    direction: this.state.joystick.dir, 
                    interval: this.state.interval
                }
            )
        );
    }

    timerTick() {
        if (this.state.joystick.force !== 0) {
            this.sendStateToSocket();
        }
    }

    updatePosition(force, angle, dir) {
        this.setState({
            joystick: {
                force: force,
                angle: angle,
                dir: dir
            }
        });

    }

    resetPosition() {
        this.setState({
            joystick: {
                force: 0,
                angle: 0,
                dir: null
            }
        });
        this.sendStateToSocket();
    }

    render() {
        let actual = this.state.joystick
        return (
            <div>
                <h1>this is some demo app</h1>
                <Joystick 
                    force={parseFloat(actual.force).toFixed(2)}
                    angle={parseFloat(actual.angle).toFixed(2)}

                    onMove={(force, angle, dir) => this.updatePosition(force, angle, dir)}
                    onStop={() => this.resetPosition()}
                />
                <Timer
                    interval={this.state.interval}
                    onUpdate={() => this.timerTick()}
                />
            </div>
        );
    }
}

export default App;
