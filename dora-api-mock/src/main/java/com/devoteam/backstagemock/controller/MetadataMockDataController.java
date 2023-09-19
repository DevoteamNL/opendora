package com.devoteam.backstagemock.controller;

import lombok.SneakyThrows;
import org.springframework.core.io.ClassPathResource;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.nio.file.Files;

@RestController
@RequestMapping("/mock-data")
public class MetadataMockDataController {

    @SneakyThrows
    @GetMapping(produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<String> getMockData() {
        var file = new ClassPathResource("stubs/home/mock-data.json").getFile();

        var data = Files.readString(file.toPath());

        return ResponseEntity.ok(data);
    }

}
