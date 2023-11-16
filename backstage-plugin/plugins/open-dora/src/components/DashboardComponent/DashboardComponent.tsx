import {
  Content,
  Header,
  Page,
  Progress,
  ResponseErrorPanel,
  SupportButton,
} from '@backstage/core-components';
import { getEntityRelations, useEntity } from '@backstage/plugin-catalog-react';
import { Grid } from '@material-ui/core';
import React from 'react';
import { MetricContext } from '../../services/MetricContext';
import { useMetricData } from '../../services/MetricDataHook';
import { BarChartComponent } from '../BarChartComponent/BarChartComponent';
import { DropdownComponent } from '../DropdownComponent/DropdownComponent';
import './DashboardComponent.css';

export interface DashboardComponentProps {
  entityName?: string;
  entityGroup?: string;
}

const ChartGridItem = ({ type, label }: { type: string; label: string }) => {
  const { chartData, error } = useMetricData(type);

  const chartOrProgressComponent = chartData ? (
    <BarChartComponent metricData={chartData} />
  ) : (
    <Progress variant="indeterminate" />
  );

  const errorOrResponse = error ? (
    <ResponseErrorPanel error={error} />
  ) : (
    chartOrProgressComponent
  );

  return (
    <Grid item xs={12} className="gridBorder">
      <div className="gridBoxText">
        <h1>{label}</h1>
        {errorOrResponse}
      </div>
    </Grid>
  );
};

export const DashboardComponent = ({
  entityName,
  entityGroup,
}: DashboardComponentProps) => {
  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('weekly');

  return (
    <MetricContext.Provider
      value={{
        aggregation: selectedTimeUnit,
        team: entityGroup,
        project: entityName,
      }}
    >
      <Page themeId="tool">
        <Header
          title="OpenDORA (by Devoteam)"
          subtitle="Through insight to perfection"
        >
          <SupportButton>Plugin for displaying DORA Metrics</SupportButton>
        </Header>
        <Content>
          <Grid container spacing={3} direction="column">
            <Grid container>
              <Grid item xs={12} className="gridBorder">
                <div className="gridBoxText">
                  <Grid container>
                    <Grid item xs={4}>
                      <DropdownComponent
                        onSelect={setSelectedTimeUnit}
                        selection={selectedTimeUnit}
                      />
                    </Grid>
                  </Grid>
                </div>
              </Grid>

              <ChartGridItem type="df_count" label="Deployment Frequency" />
              <ChartGridItem
                type="df_average"
                label="Deployment Frequency Average"
              />
            </Grid>
            <Grid item />
          </Grid>
        </Content>
      </Page>
    </MetricContext.Provider>
  );
};

export const EntityDashboardComponent = () => {
  const { entity } = useEntity();
  const entityGroup = getEntityRelations(entity, 'ownedBy')[0]?.name;
  const entityName = entity.metadata.name;

  return (
    <DashboardComponent entityName={entityName} entityGroup={entityGroup} />
  );
};
