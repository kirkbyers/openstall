import { h, Component } from 'preact';
import { icon } from '@fortawesome/fontawesome-svg-core'
import { faDoorOpen, faDoorClosed } from '@fortawesome/free-solid-svg-icons'

function updateMonitor(newMonitor, monitorArray) {
    const result = monitorArray.slice(0);
    let found = -1;
    for (let i = 0; i < monitorArray.length; i++) {
        if (newMonitor.id == monitorArray[i].id) {
            result[i] = newMonitor;
            found = i;
            break
        }
    }
    if (found < 0) {
        result.push(newMonitor)
    }
    return result
}

function getMonitors() {
    return fetch(`//${process.env.API_BASE}/status`)
}

export class DoorMonitors extends Component {
    constructor() {
        super();
        this.state = {
            monitors: [],
            wsConn: new WebSocket(`ws://${process.env.API_BASE}/sub`),
            wsConnected: true
        };
        this.state.wsConn.onclose = (evt) => {
            this.setState(() => ({ wsConnected: false }));
        };
        this.state.wsConn.onmessage = (evt) => {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                const parsedMonitor = JSON.parse(messages[i]);
                this.setState(() => ({ monitors: updateMonitor(parsedMonitor, this.state.monitors) }));
            }
        }
        getMonitors().then((res) => {
            return res.json()
        }).then((resJson) => {
            this.setState(() => ({ monitors: resJson }))
        });
    }
    render(props, state) {
        return (
            <div>
                {state.monitors.map(((monitor) => {
                    return (
                        <DoorMonitor monitor={monitor} />
                    )
                }))}
            </div>
        )
    }
}

export const DoorMonitor = (props) => {
    const iconConfig = {
        transform: {
        },
        styles: { 'background-color': 'white', 'color': 'blue' },
    };
    const { name, status } = props.monitor;
    const openDoorIcon = icon(faDoorOpen, iconConfig);
    const closedDoorIcon = icon(faDoorClosed, iconConfig);
    const activeIcon = status === 'open' ? openDoorIcon : closedDoorIcon;
    return (
        <div>
            <h3>
                {name}
            </h3>
            <div dangerouslySetInnerHTML={{ __html: activeIcon.html }} />
            <p>
                is <b><i>{status}</i></b>.
            </p>
        </div>
    )
}