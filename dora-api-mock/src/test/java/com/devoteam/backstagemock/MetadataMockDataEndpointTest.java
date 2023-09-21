package com.devoteam.backstagemock;

import lombok.SneakyThrows;
import org.hamcrest.Matchers;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.client.MockRestServiceServer;
import org.springframework.web.client.RestTemplate;

import static org.springframework.test.web.client.ExpectedCount.once;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.*;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class MetadataMockDataEndpointTest {

    @Test
    @SneakyThrows
    void mockDataEndpointShouldReturnMetadata() {
        var restTemplate = new RestTemplate();
        var server = MockRestServiceServer.bindTo(restTemplate).build();

        var metadata = Matchers.allOf(
                Matchers.hasKey("id"),
                Matchers.hasKey("iid"),
                Matchers.hasKey("project_id"),
                Matchers.hasKey("status"),
                Matchers.hasKey("ref"),
                Matchers.hasKey("sha"),
                Matchers.hasKey("before_sha"),
                Matchers.hasKey("tag"),
                Matchers.hasKey("yaml_errors"),
                Matchers.hasKey("user"),
                Matchers.hasKey("created_at"),
                Matchers.hasKey("updated_at"),
                Matchers.hasKey("started_at"),
                Matchers.hasKey("finished_at"),
                Matchers.hasKey("committed_at"),
                Matchers.hasKey("duration"),
                Matchers.hasKey("queued_duration"),
                Matchers.hasKey("coverage"),
                Matchers.hasKey("web_url")
        );

        var user = Matchers.allOf(
                Matchers.hasKey("id"),
                Matchers.hasKey("name"),
                Matchers.hasKey("username"),
                Matchers.hasKey("state"),
                Matchers.hasKey("avatar_url"),
                Matchers.hasKey("web_url")
        );

        server.expect(once(), requestTo("http://localhost:10666/mock-data"))
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.metadata_1").exists())
                .andExpect(jsonPath("$.metadata_1").value(metadata))
                .andExpect(jsonPath("$.metadata_1.user").value(user))
                .andExpect(jsonPath("$.metadata_2").exists())
                .andExpect(jsonPath("$.metadata_2").value(metadata))
                .andExpect(jsonPath("$.metadata_2.user").value(user))
                .andExpect(jsonPath("$.metadata_3").exists())
                .andExpect(jsonPath("$.metadata_3").value(metadata))
                .andExpect(jsonPath("$.metadata_3.user").value(user));
    }

}
