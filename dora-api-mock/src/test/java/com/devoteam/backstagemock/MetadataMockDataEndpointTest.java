package com.devoteam.backstagemock;

import com.jayway.jsonpath.Configuration;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;

import static com.jayway.jsonpath.matchers.JsonPathMatchers.hasJsonPath;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class MetadataMockDataEndpointTest {

    @Autowired
    private TestRestTemplate restTemplate;

    @Test
    void mockDataEndpointShouldReturnMetadata() {
        var metadataMatcher = allOf(
                hasKey("id"),
                hasKey("iid"),
                hasKey("project_id"),
                hasKey("status"),
                hasKey("ref"),
                hasKey("sha"),
                hasKey("before_sha"),
                hasKey("tag"),
                hasKey("yaml_errors"),
                hasKey("user"),
                hasKey("created_at"),
                hasKey("updated_at"),
                hasKey("started_at"),
                hasKey("finished_at"),
                hasKey("committed_at"),
                hasKey("duration"),
                hasKey("queued_duration"),
                hasKey("coverage"),
                hasKey("web_url")
        );

        var userMatcher = allOf(
                hasKey("id"),
                hasKey("name"),
                hasKey("username"),
                hasKey("state"),
                hasKey("avatar_url"),
                hasKey("web_url")
        );

        var result = restTemplate.getForObject("http://localhost:10666/mock-data", String.class);

        var json = Configuration.defaultConfiguration().jsonProvider().parse(result);
        assertThat(json, hasJsonPath("$.metadata_1", is(metadataMatcher)));
        assertThat(json, hasJsonPath("$.metadata_1.user", is(userMatcher)));

        assertThat(json, hasJsonPath("$.metadata_2", is(metadataMatcher)));
        assertThat(json, hasJsonPath("$.metadata_2.user", is(userMatcher)));

        assertThat(json, hasJsonPath("$.metadata_3", is(metadataMatcher)));
        assertThat(json, hasJsonPath("$.metadata_3.user", is(userMatcher)));
    }

}
