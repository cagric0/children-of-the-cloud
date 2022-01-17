package com.coc.childrenofcloud.network;

import java.io.IOException;
import java.net.CookieManager;
import java.net.CookiePolicy;

import okhttp3.Interceptor;
import okhttp3.JavaNetCookieJar;
import okhttp3.OkHttpClient;
import okhttp3.Response;
import okhttp3.logging.HttpLoggingInterceptor;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;

public class Service {

    public static ServiceAPI serviceAPI;
//    public final static String BASE_URL = "http://192.168.1.92:3000/";
    public final static String BASE_URL = "https://children-of-the-cloud-r7ykd7ggdq-uc.a.run.app";

    public static void InitNetworking() {
        Retrofit retrofit = new Retrofit.Builder()
                .baseUrl(BASE_URL)
                .addConverterFactory(GsonConverterFactory.create())
                .client(okHttpClient())
                .build();

        serviceAPI = retrofit.create(ServiceAPI.class);
    }

    private static OkHttpClient okHttpClient() {
        HttpLoggingInterceptor interceptor = new HttpLoggingInterceptor();
        interceptor.setLevel(HttpLoggingInterceptor.Level.BODY);

        CookieManager cookieManager = new CookieManager();
        cookieManager.setCookiePolicy(CookiePolicy.ACCEPT_ALL);

        OkHttpClient.Builder builder = new OkHttpClient.Builder();
        builder.addInterceptor(new Interceptor() {
            @Override
            public Response intercept(Chain chain) throws IOException {
                Response response = chain.proceed(chain.request());
                // Do anything with response here
                //if we ant to grab a specific cookie or something..


                return response;
            }
        });

        builder.cookieJar(new JavaNetCookieJar(cookieManager));
        builder.addInterceptor(interceptor);

        return builder.build();

    }
}
