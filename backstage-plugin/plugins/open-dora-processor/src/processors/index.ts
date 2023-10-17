import { CatalogProcessor } from '@backstage/plugin-catalog-node';
import { Entity } from '@backstage/catalog-model';
import { Config } from '@backstage/config';

interface DoraAnnotations {
  projectUrl: string;
  projectId: number;
  deploymentPattern: string;
  productionPattern?: string;
}

// A processor that sets up DevLake scopes for GitLab components
export class OpenDoraDevLakeProcessor implements CatalogProcessor {
  constructor(private config: Config) {}

  getProcessorName(): string {
    return 'OpenDORA DevLake Processor';
  }

  async postProcessEntity(entity: Entity): Promise<Entity> {
    const parsedAnnotations = this.parseDoraAnnotations(entity);
    if (parsedAnnotations) {
      await this.setupDora(parsedAnnotations);
    }

    return entity;
  }

  private parseDoraAnnotations(entity: Entity): DoraAnnotations | undefined {
    const annotations = entity.metadata.annotations;
    if (
      annotations?.['devoteam.com/opendora-enabled'] !== 'true' ||
      entity.kind !== 'Component'
    ) {
      return undefined;
    }

    const {
      'devoteam.com/opendora-project-url': projectUrl,
      'devoteam.com/opendora-deployment-pattern': deploymentPattern,
      'devoteam.com/opendora-production-pattern': productionPattern,
    } = annotations;

    const projectId = Number.parseInt(
      annotations['devoteam.com/opendora-project-id'],
    );

    if (isNaN(projectId) || !projectUrl || !deploymentPattern) {
      return undefined;
    }

    return {
      projectUrl: projectUrl,
      projectId: projectId,
      deploymentPattern: deploymentPattern,
      productionPattern: productionPattern,
    };
  }

  private async setupDora(annotations: DoraAnnotations) {
    const baseUrl = this.config.getString('open-dora.devLakeBaseUrl');
    const connectionId = this.config.getNumber('open-dora.devLakeConnectionId');
    const basicAuth = this.config.getString('open-dora.devLakeBasicAuth');

    const url = new URL(baseUrl);
    url.pathname = `api/plugins/gitlab/connections/${connectionId}/scopes`;

    await fetch(url, {
      method: 'PUT',
      body: JSON.stringify({
        data: [
          {
            HttpUrlToRepo: annotations.projectUrl + '.git',
            connectionId: connectionId,
            gitlabId: annotations.projectId,
            name: annotations.projectUrl,
          },
        ],
      }),
      headers: {
        Authorization: basicAuth,
      },
    });
  }
}
