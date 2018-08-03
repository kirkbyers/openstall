import { h } from 'preact';
import { icon } from '@fortawesome/fontawesome-svg-core'
import { faDoorOpen, faDoorClosed } from '@fortawesome/free-solid-svg-icons'

export const DoorMonitor = (props) => {
    const iconConfig = {
        transform: {
            // size: 50
        },
        styles: { 'background-color': 'white', 'color': 'blue' },
        // classes: ['fa-5x']
    };

    const openDoorIcon = icon(faDoorOpen, iconConfig);
    const closedDoorIcon = icon({ iconName: 'door-closed', transform: iconConfig });
    return (
        <div>
            This is a door <span dangerouslySetInnerHTML={{ __html: openDoorIcon.html }} />
        </div>
    )
}