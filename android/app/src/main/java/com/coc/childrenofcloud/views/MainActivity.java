package com.coc.childrenofcloud.views;

import android.Manifest;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.database.Cursor;
import android.graphics.BitmapFactory;
import android.media.MediaRecorder;
import android.net.Uri;
import android.os.Bundle;
import android.os.Environment;
import android.provider.MediaStore;
import android.util.Log;
import android.view.View;
import android.widget.ImageView;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;
import androidx.core.app.ActivityCompat;
import androidx.core.content.FileProvider;
import androidx.fragment.app.FragmentActivity;
import androidx.recyclerview.widget.RecyclerView;

import com.coc.childrenofcloud.R;
import com.coc.childrenofcloud.Utils;
import com.coc.childrenofcloud.models.DetectedObject;
import com.coc.childrenofcloud.models.ResponseModel;
import com.coc.childrenofcloud.network.Service;
import com.coc.childrenofcloud.views.adapters.ContentAdapter;
import com.coc.childrenofcloud.views.adapters.GuessedAdapter;
import com.coc.childrenofcloud.views.customviews.CustomRecordButton;
import com.github.squti.androidwaverecorder.WaveRecorder;

import java.io.File;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;

import butterknife.BindView;
import butterknife.ButterKnife;
import butterknife.OnClick;
import okhttp3.MediaType;
import okhttp3.MultipartBody;
import okhttp3.RequestBody;
import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class MainActivity extends FragmentActivity {

    @BindView(R.id.cl_select)
    View selectScreen;
    @BindView(R.id.cl_selected)
    View selectedScreen;
    @BindView(R.id.cl_results)
    View resultsScreen;
    @BindView(R.id.cl_mic)
    View micScreen;
    @BindView(R.id.image)
    ImageView imageView;
    @BindView(R.id.btn_record)
    CustomRecordButton btnRecord;
    @BindView(R.id.rv_content)
    RecyclerView rvContent;
    @BindView(R.id.rv_guessed)
    RecyclerView rvGuessed;
    @BindView(R.id.tv_content_not_found)
    View notFoundContent;
    @BindView(R.id.tv_guessed_not_found)
    View notFoundGuessed;
    @BindView(R.id.iv_status)
    ImageView ivResult;

    private static final int REQUEST_RECORD_AUDIO_PERMISSION = 200;
    private boolean permissionToRecordAccepted = false;
    private String[] permissions = {Manifest.permission.RECORD_AUDIO, android.Manifest.permission.WRITE_EXTERNAL_STORAGE};

    private WaveRecorder waveRecorder;

    private File photo;
    private static String audioFileName = null;
    private String imageFromCamera;

    private ContentAdapter contentAdapter;
    private GuessedAdapter guessedAdapter;

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        switch (requestCode) {
            case REQUEST_RECORD_AUDIO_PERMISSION:
                permissionToRecordAccepted = grantResults[0] == PackageManager.PERMISSION_GRANTED;
                break;
        }
        if (!permissionToRecordAccepted) finish();

    }

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        ButterKnife.bind(this);
        btnRecord.SetUpActions(this::StartRecording, this::StopRecording);

        audioFileName = getExternalCacheDir().getAbsolutePath();
        audioFileName += "/audiorecord.wav";
        waveRecorder = new WaveRecorder(audioFileName);

        contentAdapter = new ContentAdapter();
        rvContent.setAdapter(contentAdapter);
        guessedAdapter = new GuessedAdapter();
        rvGuessed.setAdapter(guessedAdapter);

        ActivityCompat.requestPermissions(this, permissions, REQUEST_RECORD_AUDIO_PERMISSION);

        ModeSelectScreen();
    }

    private void ModeSelectScreen() {
        selectScreen.setVisibility(View.VISIBLE);
        selectedScreen.setVisibility(View.GONE);
        photo = null;
        imageFromCamera = null;
    }

    private void ModeRecord() {
        selectScreen.setVisibility(View.GONE);
        selectedScreen.setVisibility(View.VISIBLE);
        micScreen.setVisibility(View.VISIBLE);
        resultsScreen.setVisibility(View.GONE);
    }

    private void ModeResult(List<DetectedObject> objects, List<String> speech, boolean isSuccess) {
        selectScreen.setVisibility(View.GONE);
        selectedScreen.setVisibility(View.VISIBLE);
        micScreen.setVisibility(View.GONE);
        resultsScreen.setVisibility(View.VISIBLE);

        if (objects == null || objects.size() == 0) {
            rvContent.setVisibility(View.INVISIBLE);
            notFoundContent.setVisibility(View.VISIBLE);
        } else {
            rvContent.setVisibility(View.VISIBLE);
            notFoundContent.setVisibility(View.GONE);
            contentAdapter.SetData(objects);
        }

        if (speech == null || speech.size() == 0) {
            rvGuessed.setVisibility(View.INVISIBLE);
            notFoundGuessed.setVisibility(View.VISIBLE);
        } else {
            rvGuessed.setVisibility(View.VISIBLE);
            notFoundGuessed.setVisibility(View.GONE);

            guessedAdapter.SetData(speech);
        }

        ivResult.setImageResource(isSuccess ? R.drawable.ic_done : R.drawable.ic_wrong);
    }

    private void SendRequest() {
        Utils.showLoading(getSupportFragmentManager());
        MultipartBody.Part filePart = MultipartBody.Part.createFormData(
                "image",
                photo.getName(),
                RequestBody.create(MediaType.parse("image/*"), photo));

        File audioFile = new File(audioFileName);
        MultipartBody.Part audioPart = MultipartBody.Part.createFormData(
                "audio",
                audioFile.getName(),
                RequestBody.create(MediaType.parse("audio/*"), audioFile));

        Service.serviceAPI.sendData(filePart, audioPart).enqueue(new Callback<ResponseModel>() {
            @Override
            public void onResponse(Call<ResponseModel> call, Response<ResponseModel> response) {
                Utils.dismissLoading();
                if (!response.isSuccessful() || response.body() == null) {
                    Utils.ShowErrorToast(MainActivity.this, response.message());
                    return;
                }
                ModeResult(response.body().objects, response.body().speechTexts, response.body().isSuccess);
            }

            @Override
            public void onFailure(Call<ResponseModel> call, Throwable t) {
                Utils.dismissLoading();
                Utils.ShowErrorToast(MainActivity.this, t.getMessage());
            }
        });
    }

    private void StartRecording() {
        Log.d("=======", "Pressed!");
        waveRecorder.startRecording();
    }

    private void StopRecording() {
        Log.d("=======", "Released!");
        waveRecorder.stopRecording();

        SendRequest();
    }

    @OnClick(R.id.btn_select_galery)
    public void SelectFromGalery() {
        Intent i = new Intent();
        i.setType("image/*");
        i.setAction(Intent.ACTION_PICK);

        startActivityForResult(Intent.createChooser(i, "Select Picture"), 123);
    }

    public void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);

        if (resultCode == RESULT_OK) {

            if (requestCode == 123) {

                Uri selectedImageUri = data.getData();
                if (null != selectedImageUri) {
                    // update the preview image in the layout
                    imageView.setImageURI(selectedImageUri);

                    String[] filePathColumn = {MediaStore.Images.Media.DATA};
                    Cursor cursor = getContentResolver().query(selectedImageUri,
                            filePathColumn, null, null, null);
                    cursor.moveToFirst();

                    int columnIndex = cursor.getColumnIndex(filePathColumn[0]);
                    String picturePath = cursor.getString(columnIndex);
                    cursor.close();
                    photo = new File(picturePath);

                }
                ModeRecord();
            } else if (requestCode == 222) {
//                Bundle extras = data.getExtras();
//                Bitmap imageBitmap = (Bitmap) extras.get("data");
//                imageView.setImageBitmap(imageBitmap);
                imageView.setImageBitmap(BitmapFactory.decodeFile(imageFromCamera));
                photo = new File(imageFromCamera);
                ModeRecord();
            }
        }
    }

    @OnClick(R.id.btn_retry)
    public void RetryClicked() {
        ModeRecord();
    }

    @OnClick(R.id.btn_take_photo)
    public void TakeAPhoto() {

        Intent takePictureIntent = new Intent(MediaStore.ACTION_IMAGE_CAPTURE);
        // Ensure that there's a camera activity to handle the intent
        if (takePictureIntent.resolveActivity(getPackageManager()) != null) {
            // Create the File where the photo should go
            File photoFile = null;
            try {
                photoFile = createImageFile();
            } catch (IOException ex) {
                // Error occurred while creating the File
            }

            // Continue only if the File was successfully created
            if (photoFile != null) {
                Uri photoURI = FileProvider.getUriForFile(this,
                        "com.example.android.fileprovider",
                        photoFile);
                takePictureIntent.putExtra(MediaStore.EXTRA_OUTPUT, photoURI);
                startActivityForResult(takePictureIntent, 222);
            }
        }
    }

    private File createImageFile() throws IOException {
        // Create an image file name

        File storageDir = getExternalFilesDir(Environment.DIRECTORY_PICTURES);
        String timeStamp = new SimpleDateFormat("yyyyMMdd_HHmmss").format(new Date());
        String imageFileName = "JPEG_" + timeStamp + "_";
        File image = File.createTempFile(
                imageFileName,  /* prefix */
                ".jpg",         /* suffix */
                storageDir      /* directory */
        );

        // Save a file: path for use with ACTION_VIEW intents
        imageFromCamera = image.getAbsolutePath();
        return image;
    }

    @OnClick(R.id.btn_cancel)
    public void CancelClicked() {
        imageView.setImageURI(null);
        ModeSelectScreen();
    }
}
