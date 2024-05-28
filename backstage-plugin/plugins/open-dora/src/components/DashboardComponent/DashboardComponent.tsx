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
import { useTheme } from '@mui/material/styles';
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
  const theme = useTheme();
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
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
            }}
          >
            <Box
              sx={{
                bgcolor: theme.palette.background.default,
                display: 'flex',
                top: 8,
                position: 'sticky',
                zIndex: 2,
              }}
            >
              <DropdownComponent
                onSelect={setSelectedTimeUnit}
                selection={selectedTimeUnit}
              />
            </Box>
            <Box
              sx={{
                top: 8,
                display: 'flex',
                gridGap: 8,
                zIndex: 1,
                marginTop: 8,
                maxHeight: 900,
              }}
            >
              <BenchmarkGridItem type="df" />
              <BenchmarkGridItem type="mltc" />
              <BenchmarkGridItem type="cfr" />
              <BenchmarkGridItem type="mttr" />
            </Box>
            <Box
              sx={{
                top: 8,
                display: 'flex',
                gridGap: 8,
                zIndex: 1,
                marginTop: 8,
                justifyContent: 'space-evenly',
              }}
            >
              <ChartGridItem
                type="df_count"
                label={t('deployment-frequency.labels.deployment_frequency')}
              />
              <ChartGridItem
                type="df_average"
                label={t(
                  'deployment-frequency.labels.deployment_frequency_average',
                )}
              />
            </Box>
            <Box
              sx={{
                top: 8,
                display: 'flex',
                gridGap: 8,
                zIndex: 1,
                marginTop: 8,
                justifyContent: 'space-evenly',
              }}
            >
              <ChartGridItem
                type="mltc"
                label={t('lead-time.labels.median_lead_time_for_changes')}
              />
              <ChartGridItem
                type="cfr"
                label={t('failure-rate.labels.change_failure_rate')}
              />
            </Box>
            <Box
              sx={{
                top: 8,
                display: 'flex',
                gridGap: 8,
                zIndex: 1,
                marginY: 1,
              }}
            >
              <ChartGridItem
                type="mttr"
                label={t('recovery-time.labels.mean_time_to_recovery')}
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
