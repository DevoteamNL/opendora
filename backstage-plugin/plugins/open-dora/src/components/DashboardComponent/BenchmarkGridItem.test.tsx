import { act } from '@testing-library/react';
import { rest } from 'msw';
import React from 'react';
import {
  delayRequest,
  renderComponentWithApis,
} from '../../../testing/component-test-utils';
import { benchmarkUrl } from '../../../testing/mswHandlers';
import '../../i18n';
import { server } from '../../setupTests';
import { BenchmarkGridItem } from './BenchmarkGridItem';

describe('BenchmarkGridItem', () => {
  function renderBenchmarkGridItem() {
    return renderComponentWithApis(<BenchmarkGridItem type="df" />);
  }

  it('should show the benchmark returned from the server', async () => {
    const { queryByText } = await renderBenchmarkGridItem();

    expect(queryByText('On-demand')).not.toBeNull();
  });

  it('should show loading indicator when waiting on the benchmark to return', async () => {
    delayRequest({ key: 'on-demand' }, benchmarkUrl);

    const { queryByText, queryByRole, findByRole } =
      await renderBenchmarkGridItem();

    expect(await findByRole('progressbar')).not.toBeNull();
    expect(queryByText('On-demand')).toBeNull();

    await act(async () => {
      jest.runAllTimers();
    });

    expect(queryByRole('progressbar')).toBeNull();
    expect(queryByText('On-demand')).not.toBeNull();
  });

  it('should show the error returned from the service', async () => {
    server.use(
      rest.get(benchmarkUrl, (_, res, ctx) => {
        return res(ctx.status(401));
      }),
    );
    const { queryByText } = await renderBenchmarkGridItem();
    expect(queryByText('Error: Unauthorized')).not.toBeNull();
  });
});
