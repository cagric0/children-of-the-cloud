package com.coc.childrenofcloud.models;

import java.util.List;

public class DetectedObject {
    public List<Coordinate> coordinates;
    public String name;
    public float score;
    public boolean detectedByUser;
}