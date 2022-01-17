package com.coc.childrenofcloud.models;

import java.util.List;

public class ResponseModel {
    public boolean isSuccess = true;
    public List<DetectedObject> objects;
    public List<String> speechTexts;
    public String errorMsg;
}


