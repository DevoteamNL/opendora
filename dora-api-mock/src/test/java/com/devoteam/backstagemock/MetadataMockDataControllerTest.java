package com.devoteam.backstagemock;

import lombok.SneakyThrows;
import org.hamcrest.Matchers;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@AutoConfigureMockMvc
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class MetadataMockDataControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Test
    @SneakyThrows
    void mockDataEndpointShouldReturnMetadata() {
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

        this.mockMvc.perform(get("/mock-data"))
                .andExpect(status().isOk())
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
