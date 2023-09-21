import React from 'react';
import { Box, MenuItem, TextField } from '@material-ui/core';

interface DropdownComponentProps {
  onSelect: (selection: string) => void;
  selection: string;
}

function DropdownComponent({
  onSelect,
  selection,
}: DropdownComponentProps) {
  return (
    <Box sx={{ display: 'flex', m: 3, flexDirection: 'column' }}>
      <TextField
        style={{ width: '100%' }}
        variant="outlined"
        value={selection}
        onChange={e => onSelect(e.target.value)}
        select
        label="Time Unit"
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
