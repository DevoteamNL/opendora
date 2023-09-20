import React, { Dispatch, SetStateAction } from 'react';
import { Box, MenuItem, TextField } from '@material-ui/core';


export interface DropdownProps {
  onSelect: Dispatch<SetStateAction<string>>,
  selection: string,
  type: string,
}

function DropdownComponent( props: DropdownProps) {
  let type = props.type;
  let selection = props.selection;
  let onSelect = props.onSelect;
  // Fetch all the groups here and list them as options
  const selectionLabel = type === 'timeUnit' ?? 'Group';
  switch (type) {
    case 'timeUnit': {
      selectionLabel = 'Time Unit';
      break;
    }
    case 'group': {
      selectionLabel = 'Group';
      break;
    }
    default: {
      selectionLabel = 'Label';
    }
  }



  return (
    <Box sx={{ display: 'flex', m: 3, flexDirection: 'column' }}>
      <TextField
        style={{ width: '100%' }}
        variant="outlined"
        value={selection}
        onChange={e => onSelect(e.target.value)}
        select
        label={selectionLabel}
      >
        <MenuItem key={1} value="Weekly">
          Weekly
        </MenuItem>
        <MenuItem key={2} value="Monthly">
          Monthly
        </MenuItem>
        <MenuItem key={3} value="Quarterly">
          Quarterly
        </MenuItem>
      </TextField>
    </Box>
  );
}
export default DropdownComponent;
