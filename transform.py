from random import random

import torch
from PIL import ImageDraw


# 画像の端の方に線を引く
class RandomLine(torch.nn.Module):
    def __init__(self, p=0.5, thickness=1):
        super(RandomLine, self).__init__()
        self.p = p
        self.thickness = thickness

    def forward(self, x):
        draw = ImageDraw.Draw(x)
        width, height = x.size
        if torch.rand(1) < self.p:
            x0 = random() * width / 5
            x1 = random() * width / 5
            y0 = random() * height / 5
            y1 = random() * height / 5 + height * 4 / 5
            draw.line((x0, y0, x1, y1), fill=0, width=self.thickness)

        if torch.rand(1) < self.p:
            x0 = random() * width / 5 + width * 4 / 5
            x1 = random() * width / 5 + width * 4 / 5
            y0 = random() * height / 5
            y1 = random() * height / 5 + height * 4 / 5
            draw.line((x0, y0, x1, y1), fill=0, width=self.thickness)

        if torch.rand(1) < self.p:
            x0 = random() * width / 5
            x1 = random() * width / 5 + width * 4 / 5
            y0 = random() * height / 5
            y1 = random() * height / 5
            draw.line((x0, y0, x1, y1), fill=0, width=self.thickness)

        if torch.rand(1) < self.p:
            x0 = random() * width / 5
            x1 = random() * width / 5 + width * 4 / 5
            y0 = random() * height / 5 + height * 4 / 5
            y1 = random() * height / 5 + height * 4 / 5
            draw.line((x0, y0, x1, y1), fill=0, width=self.thickness)
        return x


if __name__ == "__main__":
    import torchvision.transforms as transforms
    from PIL import Image

    img = Image.open("data/train/train/1/1.png")
    trans = transforms.Compose([RandomLine(p=0.2)])
    trans(img).show()
