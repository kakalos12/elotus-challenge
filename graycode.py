class Solution:
    def grayCode(self, n: int) -> list[int]:
        result = [0]

        for i in range(n):
            adder = 1 << i
            for val in reversed(result):
                result.append(val + adder)

        return result
