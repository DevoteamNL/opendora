import React from 'react';
import { HighlightTextBoxComponent } from './HighlightTextBoxComponent';
import { render } from '@testing-library/react';

describe('HighlightTextBoxComponent', () => {
  it('should create a highlighted text label', async () => {
    const { getByText, queryByText } = render(
      <HighlightTextBoxComponent
        title="example title"
        healthStatus="low"
        highlight="example highlight"
      />,
    );

    expect(queryByText('example title')).not.toBeNull();

    const highlightElement = getByText('example highlight');
    expect(highlightElement.parentElement).toHaveClass('low');

    expect(queryByText('example text')).toBeNull();
  });
  it('should create a label with a notification', async () => {
    const { queryByText } = render(
      <HighlightTextBoxComponent
        title="example title"
        healthStatus="low"
        highlight="example highlight"
        text="example text"
      />,
    );

    expect(queryByText('example text')).not.toBeNull();
  });
});
