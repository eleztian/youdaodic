package main

import (
	"eleztian/youdaoDic/youdao"
	"fmt"
)

func main() {
	youdao.Config("2ccac2276928012f", "tm0aC9BVe2qZq4DuHRR9p5KdEA7y6l1Y")
	fmt.Println(youdao.Translate(
		`First Flight
　　Mr. Johnson had never been up in an aerophane before and he had read a lot about air accidents, so one day when a friend offered to take him for a ride in his own small phane, Mr. Johnson was very worried about accepting. Finally, however, his friend persuaded him that it was very safe, and Mr. Johnson boarded the plane.
　　His friend started the engine and began to taxi onto the runway of the airport. Mr. Johnson had heard that the most dangerous part of a flight were the take-off and the landing, so he was extremely frightened and closed his eyes.
　　After a minute or two he opened them again, looked out of the window of the plane, and said to his friend, "Look at those people down there. They look as small as ants, don't they?"
　　"Those are ants," answered his friend. "We're still on the ground."`,
		youdao.English, youdao.Chinese, true))
}
