package com.coc.childrenofcloud.app;

import android.content.Context;

import androidx.appcompat.app.AppCompatDelegate;

import com.coc.childrenofcloud.network.Service;


public class Application extends android.app.Application {

    private static Context context;



    public static Context getAppContext() {
        return context;
    }

    @Override
    public void onCreate() {
        super.onCreate();
        AppCompatDelegate.setDefaultNightMode(AppCompatDelegate.MODE_NIGHT_NO);
        Application.context = getApplicationContext();
        Service.InitNetworking();
    }
}
