package com.devoteam.backstagemock;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.ConfigurationPropertiesScan;

@SpringBootApplication
@ConfigurationPropertiesScan("com.devoteam.backstagemock.config")
public class DoraMockApiApplication {

    public static void main(String[] args) {
        SpringApplication.run(DoraMockApiApplication.class, args);
    }

}
