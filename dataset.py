from glob import glob
import os

import torch
from torch.utils.data import Dataset
from torchvision import transforms
from PIL import Image


class OCRDataset(Dataset):
    def __init__(self, mode: str, transform=None):
        self.data = []
        if mode == "train":
            dirs = glob("data/train/train/**")
        elif mode == "valid":
            dirs = glob("data/train/valid/**")
        elif mode == "test":
            dirs = glob("data/train/test/**")
        else:
            raise ValueError("mode must be one of 'train', 'valid', or 'test'")

        for d in dirs:
            label = os.path.basename(d)
            files = glob(d + "/*.png")
            for f in files:
                self.data.append({"label": label, "path": f})
        self.transform = transform

    def __len__(self) -> int:
        return len(self.data)

    def __getitem__(self, idx) -> (torch.Tensor, int):
        item = self.data[idx]
        img = Image.open(item["path"])
        if self.transform:
            img = self.transform(img)
        # label to one-hot encoding
        # label is between 0 to 9
        label = torch.zeros(10)
        label[int(item["label"])] = 1
        return img, label


# test
if __name__ == "__main__":
    ds = OCRDataset(transforms.ToTensor())
    print(len(ds))
    img, label = ds[0]
    print(label)
    print(img.shape)
