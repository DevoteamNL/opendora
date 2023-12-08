import {
  Content,
  Header,
  Page,
  SupportButton,
} from '@backstage/core-components';
import { getEntityRelations, useEntity } from '@backstage/plugin-catalog-react';
import { Grid } from '@material-ui/core';
import React from 'react';
import { useTranslation } from 'react-i18next';
import '../../i18n';
import { MetricContext } from '../../services/MetricContext';
import { DropdownComponent } from '../DropdownComponent/DropdownComponent';
import { BenchmarkGridItem } from './BenchmarkGridItem';
import { ChartGridItem } from './ChartGridItem';
import './DashboardComponent.css';

export interface DashboardComponentProps {
  entityName?: string;
  entityGroup?: string;
}

export const DashboardComponent = ({
  entityName,
  entityGroup,
}: DashboardComponentProps) => {
  const [t] = useTranslation();
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
              <BenchmarkGridItem type="df" />
              <ChartGridItem
                type="df_count"
                label={t('deployment_frequency.labels.deployment_frequency')}
              />
              <ChartGridItem
                type="df_average"
                label={t(
                  'deployment_frequency.labels.deployment_frequency_average',
                )}
              />
              <ChartGridItem
                type="mltc"
                label={t('lead-time.labels.median_lead_time_for_changes')}
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
