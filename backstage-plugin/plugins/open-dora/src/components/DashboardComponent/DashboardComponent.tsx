import {
  Content,
  Header,
  Page,
  SupportButton,
} from '@backstage/core-components';
import { getEntityRelations, useEntity } from '@backstage/plugin-catalog-react';
import { Box } from '@material-ui/core';
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
          <Box sx={{ display: 'flex', flexDirection: 'column' }}>
            <Box
              sx={{
                bgcolor: '#424242',
                position: 'sticky',
                top: 8,
                display: 'flex',
                gridGap: 8,
                zIndex: 1,
                boxShadow:
                  '0px 2px 4px -1px rgba(0,0,0,0.2),0px 4px 5px 0px rgba(0,0,0,0.14),0px 1px 10px 0px rgba(0,0,0,0.12)',
              }}
            >
              <DropdownComponent
                onSelect={setSelectedTimeUnit}
                selection={selectedTimeUnit}
              />
            </Box>
            <Box
              sx={{
                display: 'flex',
                marginTop: 8,
                gridGap: 8,
                maxHeight: 180,
              }}
            >
              <BenchmarkGridItem type="df" />
              <Box sx={{ flex: 3 }}>
                {/* Placeholder for other benchmark items */}
              </Box>
            </Box>
            <Box
              sx={{
                display: 'flex',
                marginTop: 8,
                gridGap: 8,
                justifyContent: 'space-evenly',
              }}
            >
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
            </Box>
            <Box sx={{ display: 'flex', marginY: 1, gridGap: 8 }}>
              <ChartGridItem
                type="mltc"
                label={t('lead-time.labels.median_lead_time_for_changes')}
              />
              <Box sx={{ flex: 1 }}>
                {/* Placeholder for other chart items */}
              </Box>
            </Box>
          </Box>
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
