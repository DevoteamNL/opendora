import React, { Props } from "react";
import { Box, MenuItem, TextField } from "@material-ui/core";


function DropdownComponent({ onSelect, selection, type }) {
// TODO

// Fetch all the groups here and list them as options
let selectionLabel = "";
    switch(type) {
        case "timeUnit": {
            selectionLabel = "Time Unit"
            break;
        }
        case "group": {
            selectionLabel = "Group"
            break;
        }
        default: {
            selectionLabel = "Label"
        }
    }

    return (
    <Box sx={{ display: "flex", m: 3, flexDirection: "column" }}>
        <TextField
          style={{ width: "100%" }}
          variant="outlined"
          value={selection}
          onChange={(e) => onSelect(e.target.value)}
          select
          label={selectionLabel}
        >
          <MenuItem key={1} value="Grp1">
            Group 1
          </MenuItem>
          <MenuItem key={2} value="Grp2">
            Group 2
          </MenuItem>
        </TextField>
    </Box>
    );

  }
  export default DropdownComponent;