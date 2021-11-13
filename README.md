# TOKEN-PRICE-MONITOR

Token được list lên PancakeSwap sẽ có contract để lấy giá tương ứng của `pair` đó.

Trong ví dụ này mình sẽ lấy giá BAKE/USDT trên PancakeSwap như sau.

Token BAKE trên bscscan.com ở đây https://bscscan.com/token/0xE02dF9e3e622DeBdD69fb838bB799E3F168902c5

Ta có được contract `0xE02dF9e3e622DeBdD69fb838bB799E3F168902c5`

Tìm contract này trên PancakeSwap https://pancakeswap.finance/info/tokens

![](https://i.imgur.com/vCAkCGd.png)

Mình sẽ lấy contract USDT/CAKE như bên dưới

![](https://i.imgur.com/N2i6vfT.png)

Contract: https://pancakeswap.finance/info/pool/0xf08046a9a44913536b6f563e33d403beb0fb3cbf => `0xf08046a9a44913536b6f563e33d403beb0fb3cbf`

Mình sẽ tìm contract này trên BscScan sẽ tìm được function `getReserves`, function này trả về:
- Reserve0: số lượng token 1 (trường hợp này là USDT)
- Reserve1: số lượng token 2 (trường hợp này là BAKE)
- BlockTimestampLast

![](https://i.imgur.com/bBkvst0.png)

Vậy giá USDT/BAKE sẽ là Reserve1/Reserve0, giá BAKE/USDT là Reserve0/Reserve1

![](https://i.imgur.com/xyj2hBD.png)

CAKE/USDT: 2.098389
USDT/CAKE: 0.476556
