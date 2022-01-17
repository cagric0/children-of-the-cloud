package com.coc.childrenofcloud.network;


import com.coc.childrenofcloud.models.ResponseModel;

import java.util.List;

import okhttp3.MultipartBody;
import retrofit2.Call;
import retrofit2.http.Body;
import retrofit2.http.Multipart;
import retrofit2.http.POST;
import retrofit2.http.Part;

public interface ServiceAPI {

    @Multipart
    @POST("final")
    Call<ResponseModel> sendData(@Part MultipartBody.Part image, @Part MultipartBody.Part audio);


}

