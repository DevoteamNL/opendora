import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { DropdownComponent } from './DropdownComponent';
import '../../i18n';

describe('DropdownComponent', () => {
  it('should create show a dropdown with the aggregation choices', async () => {
    const { getByText, queryByText } = render(
      <DropdownComponent selection="weekly" onSelect={jest.fn()} />,
    );

    expect(queryByText('Weekly')).not.toBeNull();

    const triggerElement = getByText('Weekly');
    fireEvent.mouseDown(triggerElement);

    // Trigger button and option item
    expect(screen.queryAllByText('Weekly')).toHaveLength(2);
    expect(screen.queryAllByText('Quarterly')).toHaveLength(1);
    expect(screen.queryAllByText('Monthly')).toHaveLength(1);
  });

  it('should show the selection by default', async () => {
    const { queryByText } = render(
      <DropdownComponent selection="monthly" onSelect={jest.fn()} />,
    );

    expect(queryByText('Monthly')).not.toBeNull();
    expect(queryByText('Weekly')).toBeNull();
  });

  it('should output the selection when an option is selected', async () => {
    const onSelectSpy = jest.fn();
    const { getByText } = render(
      <DropdownComponent selection="weekly" onSelect={onSelectSpy} />,
    );

    expect(onSelectSpy).not.toHaveBeenCalled();
    const triggerElement = getByText('Weekly');
    fireEvent.mouseDown(triggerElement);

    expect(onSelectSpy).not.toHaveBeenCalled();
    const monthlyElement = screen.getByText('Monthly');
    fireEvent.click(monthlyElement);

    expect(onSelectSpy).toHaveBeenCalledWith('monthly');
    expect(onSelectSpy).toHaveBeenCalledTimes(1);

    const newTriggerElement = getByText('Monthly');
    fireEvent.mouseDown(newTriggerElement);
    const quarterlyElement = screen.getByText('Quarterly');
    fireEvent.click(quarterlyElement);

    expect(onSelectSpy).toHaveBeenLastCalledWith('quarterly');
    expect(onSelectSpy).toHaveBeenCalledTimes(2);
  });
});
