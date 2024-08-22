package controller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type AnimeInfo struct {
	NameEn    string
	NameZh    string
	NameJp    string
	Season    int
	SeasonRaw string
	Episode   int
	Sub       string
	Dpi       string
	Source    string
	Group     string
}

var (
	RESOLUTION_RE = regexp.MustCompile(`1080|720|2160|4K`)
	SOURCE_RE     = regexp.MustCompile(`B-Global|[Bb]aha|[Bb]ilibili|AT-X|Web`)
	SUB_RE        = regexp.MustCompile(`[简繁日字幕]|CH|BIG5|GB`)
	TITLE_RE      = regexp.MustCompile(`(.*|\[.*])( -? \d+|\[\d+]|\[\d+.?[vV]\d]|第\d+[话話集]|\[第?\d+[话話集]]|\[\d+.?END]|[Ee][Pp]?\d+)(.*)`)
	EPISODE_RE    = regexp.MustCompile(`\d+`)
	PREFIX_RE     = regexp.MustCompile(`[^\w\s\p{Han}\p{Hiragana}\p{Katakana}-]`)
)

var CHINESE_NUMBER_MAP = map[string]int{
	"一": 1,
	"二": 2,
	"三": 3,
	"四": 4,
	"五": 5,
	"六": 6,
	"七": 7,
	"八": 8,
	"九": 9,
	"十": 10,
}

func seasonProcess(seasonInfo string) (string, string, int) {
	nameSeason := seasonInfo
	// 去除「新番」信息, 若有
	// nameSeason = regexp.MustCompile(".*新番.").ReplaceAllString(seasonInfo, "")
	nameSeason = regexp.MustCompile(`^[^\[\]】]*[\[\]】]`).ReplaceAllString(nameSeason, "")
	nameSeason = strings.TrimSpace(nameSeason)

	seasonRule := `S\d{1,2}|Season \d{1,2}|第.{1,2}[季期]|部分`
	nameSeason = regexp.MustCompile(`[\[\]]`).ReplaceAllString(nameSeason, " ")
	seasons := regexp.MustCompile(seasonRule).FindAllString(nameSeason, -1)

	if len(seasons) == 0 {
		return nameSeason, "", 1
	}

	name := regexp.MustCompile(seasonRule).ReplaceAllString(nameSeason, "")
	var season int
	var seasonRaw string

	seasonRegxp := regexp.MustCompile(`Season|S`)
	seasonPartRegxp := regexp.MustCompile(`第.{1,2}[季期]|部分`)
	for _, s := range seasons {
		seasonRaw = s
		if matched := seasonRegxp.MatchString(s); matched {
			season, _ = strconv.Atoi(regexp.MustCompile(`Season|S`).ReplaceAllString(s, ""))
			break
		} else if matched := seasonPartRegxp.MatchString(s); matched {
			seasonPro := regexp.MustCompile(`[第季期 ]`).ReplaceAllString(s, "")
			if num, err := strconv.Atoi(seasonPro); err == nil {
				season = num
			} else {
				season = CHINESE_NUMBER_MAP[seasonPro]
			}
			break
		}
	}

	return name, seasonRaw, season
}

func Parse(raw string) (info AnimeInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	info = process(raw)

	return info, nil
}

func nameProcess(name string) (string, string, string) {
	var nameEn, nameZh, nameJp string
	name = strings.TrimSpace(name)
	name = regexp.MustCompile(`[(（]仅限港澳台地区[）)]`).ReplaceAllString(name, "")
	split := regexp.MustCompile(`/|\s{2}|-\s{2}`).Split(name, -1)

	split = removeEmptyStrings(split)

	if len(split) == 1 {
		if strings.Contains(name, "_") {
			split = strings.Split(name, "_")
		} else if strings.Contains(name, " - ") {
			split = strings.Split(name, " - ")
		}
	}

	if len(split) == 1 {
		splitSpace := strings.Split(split[0], " ")
		for _, idx := range []int{0, len(splitSpace) - 1} {
			if idx < len(splitSpace) && regexp.MustCompile(`^[\p{Han}]{2,}`).MatchString(splitSpace[idx]) {
				chs := splitSpace[idx]
				splitSpace = append(splitSpace[:idx], splitSpace[idx+1:]...)
				split = []string{chs, strings.Join(splitSpace, " ")}
				break
			}
		}
	}

	for _, item := range split {
		if regexp.MustCompile(`[\p{Hiragana}\p{Katakana}]{2,}`).MatchString(item) && nameJp == "" {
			nameJp = strings.TrimSpace(item)
		} else if regexp.MustCompile(`[\p{Han}]{2,}`).MatchString(item) && nameZh == "" {
			nameZh = strings.TrimSpace(item)
		} else if regexp.MustCompile(`[a-zA-Z]{3,}`).MatchString(item) && nameEn == "" {
			nameEn = strings.TrimSpace(item)
		}
	}

	return nameEn, nameZh, nameJp
}

func process(rawTitle string) AnimeInfo {
	processedTitle := strings.TrimSpace(strings.ReplaceAll(rawTitle, "\n", ""))
	processedTitle = preProcess(processedTitle)
	group := getGroup(processedTitle)
	season_info, episode_info, other := getTheTree(processedTitle)
	prefix := prefixProcess(season_info, group)
	raw_name, season_raw, season := seasonProcess(prefix)
	name_en, name_zh, name_jp := nameProcess(raw_name)
	rawEpisode := EPISODE_RE.FindString(episode_info)
	episode, _ := strconv.Atoi(rawEpisode)
	sub, resolution, source := findTags(other)

	return AnimeInfo{
		NameEn:    name_en,
		NameZh:    name_zh,
		NameJp:    name_jp,
		Season:    season,
		SeasonRaw: season_raw,
		Episode:   episode,
		Sub:       sub,
		Dpi:       resolution,
		Source:    source,
		Group:     group,
	}
}

func prefixProcess(raw string, group string) string {
	// 替换 .{group}.
	re := regexp.MustCompile(fmt.Sprintf(".%s.", group))
	raw = re.ReplaceAllString(raw, "")

	// 替换 PREFIX_RE.sub("/", raw)
	rawProcess := PREFIX_RE.ReplaceAllString(raw, "/")

	argGroup := strings.Split(rawProcess, "/")
	argGroup = removeEmptyStrings(argGroup)

	if len(argGroup) == 1 {
		argGroup = strings.Split(argGroup[0], " ")
	}
	Xinfan := regexp.MustCompile(`新番|月?番`)
	GanAoTai := regexp.MustCompile(`港澳台地区`)
	for _, arg := range argGroup {
		if match := Xinfan.MatchString(arg); match && len(arg) <= 5 {
			re = regexp.MustCompile(fmt.Sprintf(".%s.", arg))
			raw = re.ReplaceAllString(raw, "")
		} else if match := GanAoTai.MatchString(arg); match {
			re = regexp.MustCompile(fmt.Sprintf(".%s.", arg))
			raw = re.ReplaceAllString(raw, "")
		}
	}

	return raw
}

func removeEmptyStrings(slice []string) []string {
	var result []string
	for _, str := range slice {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func getTheTree(precessedTitle string) (string, string, string) {
	match := TITLE_RE.FindStringSubmatch(precessedTitle)
	return strings.TrimSpace(match[1]), strings.TrimSpace(match[2]), strings.TrimSpace(match[3])
}

func preProcess(rawName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(rawName, "【", "["), "】", "]")

}

func getGroup(name string) string {
	re := regexp.MustCompile(`[\[\]]`)
	parts := re.Split(name, -1)
	if len(parts) > 1 {

		return parts[1]
	}
	return ""
}

func findTags(other string) (string, string, string) {
	CLEAN_RE := regexp.MustCompile(`[\[\]()（）]`)

	cleanedOther := CLEAN_RE.ReplaceAllString(other, " ")
	elements := strings.Fields(cleanedOther)

	var sub, resolution, source string

	for _, element := range elements {
		if element == "" {
			continue
		}
		if SUB_RE.MatchString(element) {
			sub = element
		} else if RESOLUTION_RE.MatchString(element) {
			resolution = element
		} else if SOURCE_RE.MatchString(element) {
			source = element
		}
	}

	sub = cleanSub(sub)

	return sub, resolution, source
}

func cleanSub(sub string) string {
	if sub == "" {
		return sub
	}
	return regexp.MustCompile(`_MP4|_MKV`).ReplaceAllString(sub, "")
}
