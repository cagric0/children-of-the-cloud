package com.coc.childrenofcloud.views.adapters;

import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageView;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.RecyclerView;

import com.coc.childrenofcloud.R;
import com.coc.childrenofcloud.models.DetectedObject;

import java.util.List;

import butterknife.BindView;
import butterknife.ButterKnife;

public class ContentAdapter extends RecyclerView.Adapter<ContentAdapter.ViewHolder> {

    private List<DetectedObject> detectedObjects;

    public void SetData(List<DetectedObject> detectedObjects) {
        this.detectedObjects = detectedObjects;
        notifyDataSetChanged();
    }

    @NonNull
    @Override
    public ViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View v = LayoutInflater.from(parent.getContext()).inflate(R.layout.item_found, parent, false);
        return new ViewHolder(v);
    }

    @Override
    public void onBindViewHolder(@NonNull ViewHolder holder, int position) {
        holder.BindData(detectedObjects.get(position));
    }

    @Override
    public int getItemCount() {
        return detectedObjects == null ? 0 : detectedObjects.size();
    }

    class ViewHolder extends RecyclerView.ViewHolder {

        @BindView(R.id.tv_item)
        TextView tvItem;
        @BindView(R.id.iv_item)
        ImageView ivItem;

        public ViewHolder(@NonNull View itemView) {
            super(itemView);
            ButterKnife.bind(this, itemView);
        }

        public void BindData(DetectedObject object) {
            tvItem.setText(object.name);
            ivItem.setImageResource(object.detectedByUser ? R.drawable.ic_done : R.drawable.ic_wrong);
        }
    }
}
