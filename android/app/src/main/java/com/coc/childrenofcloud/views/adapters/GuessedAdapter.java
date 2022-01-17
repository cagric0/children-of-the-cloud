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

public class GuessedAdapter extends RecyclerView.Adapter<GuessedAdapter.ViewHolder> {

    private List<String> detectedObjects;

    public void SetData(List<String> detectedObjects) {
        this.detectedObjects = detectedObjects;
        notifyDataSetChanged();
    }

    @NonNull
    @Override
    public ViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View v = LayoutInflater.from(parent.getContext()).inflate(R.layout.item_guessed, parent, false);
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

        public ViewHolder(@NonNull View itemView) {
            super(itemView);
            ButterKnife.bind(this, itemView);
        }

        public void BindData(String object) {
            tvItem.setText(object);
        }
    }
}
