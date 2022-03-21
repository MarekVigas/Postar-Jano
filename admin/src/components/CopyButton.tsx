import {OverlayTrigger, Popover} from "react-bootstrap";
import React from "react";

interface Props {
   selector :string;
}

const CopyButton :React.FC<Props> = ({selector}) => {
    const popover = (
        <Popover id="popover-basic">
            <Popover.Title as="h3">Skopirovane!</Popover.Title>
            <Popover.Content>Pouzi ctrl+v</Popover.Content>
        </Popover>
    );

    const copyTable = () => {
        const elTable = document.querySelector(selector);

        let range, sel;

        // Ensure that range and selection are supported by the browsers
        if (document.createRange && window.getSelection) {

            range = document.createRange();
            sel = window.getSelection();
            // unselect any element in the page
            if (sel != null) {
                sel.removeAllRanges();
            }

            if (sel == null) return

            if (elTable == null) return
            try {
                range.selectNodeContents(elTable);
                sel.addRange(range);
            } catch (e) {
                range.selectNode(elTable);
                sel.addRange(range);
            }

            document.execCommand('copy');
            sel.removeAllRanges();
        }
    }

    return (
        <OverlayTrigger trigger="click" placement="right" overlay={popover} >
            <button onClick={() => copyTable()}>Kopirovat</button>
        </OverlayTrigger>
    )
}

export default CopyButton