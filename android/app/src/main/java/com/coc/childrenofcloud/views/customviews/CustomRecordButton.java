package com.coc.childrenofcloud.views.customviews;

import android.content.Context;
import android.util.AttributeSet;
import android.view.MotionEvent;
import android.view.View;
import android.widget.Button;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;


public class CustomRecordButton extends androidx.appcompat.widget.AppCompatButton {

    private IAction onDownAction;
    private IAction onUpAction;

    public CustomRecordButton(@NonNull Context context) {
        super(context);
    }

    public CustomRecordButton(@NonNull Context context, @Nullable AttributeSet attrs) {
        super(context, attrs);
    }

    public CustomRecordButton(@NonNull Context context, @Nullable AttributeSet attrs, int defStyleAttr) {
        super(context, attrs, defStyleAttr);
    }

    public void SetUpActions(IAction onDownAction, IAction onUpAction) {
        this.onDownAction = onDownAction;
        this.onUpAction = onUpAction;
    }

    @Override
    public boolean onTouchEvent(@NonNull MotionEvent ev) {

        if (ev.getAction() == MotionEvent.ACTION_DOWN) {
            if (onDownAction != null)
                onDownAction.OnAction();
            setPressed(true);
            return true;
        } else if (ev.getAction() == MotionEvent.ACTION_UP) {
            if (onUpAction != null)
                onUpAction.OnAction();
            setPressed(false);
            return true;
        }
        return true;
    }
}
