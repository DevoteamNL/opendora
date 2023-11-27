import { Box, MenuItem, TextField } from '@material-ui/core';
import React from 'react';

import { useTranslation } from 'react-i18next';

interface DropdownComponentProps {
  onSelect: (selection: string) => void;
  selection: string;
}

export const DropdownComponent = ({
  onSelect,
  selection,
}: DropdownComponentProps) => {
  const [t] = useTranslation();
  return (
    <Box sx={{ display: 'flex', m: 3, flexDirection: 'column' }}>
      <TextField
        style={{ width: '100%' }}
        variant="outlined"
        value={selection}
        onChange={e => onSelect(e.target.value)}
        select
        label={t('dropdown_time_units.time_unit')}
      >
        <MenuItem key={1} value="weekly">
          {t('dropdown_time_units.weekly')}
        </MenuItem>
        <MenuItem key={2} value="monthly">
          {t('dropdown_time_units.monthly')}
        </MenuItem>
        <MenuItem key={3} value="quarterly">
          {t('dropdown_time_units.quarterly')}
        </MenuItem>
      </TextField>
    </Box>
  );
};
