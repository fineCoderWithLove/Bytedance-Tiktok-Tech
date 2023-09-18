package util

import (
	"bufio"
	"demotest/douyin-api/globalinit/constant"
	"go.uber.org/zap"
	"github.com/mozillazg/go-pinyin"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// TrieNode 表示字典树的节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// WordsTrie 表示敏感词字典树
type WordsTrie struct {
	root *TrieNode
}

// SearchResult 表示匹配结果
type SearchResult struct {
	UserID      int64
	CurrentTime time.Time
	Text        string
	Words       []string
}

// NewTrieNode 创建一个字典树节点
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
}

var WorldFilter *WordsTrie

// NewWordsTrieFromFile 从文件中加载敏感词并创建字典树
func NewWordsTrieFromFile(filePath string) (*WordsTrie, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	trie := &WordsTrie{root: NewTrieNode()}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		pinyinSlice := pinyin.Pinyin(word, pinyin.NewArgs())
		if len(pinyinSlice) != 0 {
			var result strings.Builder
			for _, row := range pinyinSlice {
				result.WriteString(strings.Join(row, ", "))
			}
			trie.insertWord(strings.ToLower(result.String()))
		}
		trie.insertWord(strings.ToLower(word)) // 将敏感词转换为小写并插入字典树
	}

	return trie, nil
}

// insertWord 向字典树中插入一个词语（小写）
func (trie *WordsTrie) insertWord(word string) {
	node := trie.root
	for _, char := range word {
		if !unicode.IsLetter(char) {
			// 跳过非字母字符
			continue
		}
		if _, ok := node.children[char]; !ok {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// FindWords 返回文本中包含的敏感词列表（忽略大小写）
func (trie *WordsTrie) FindWords(text string, userID int64) SearchResult {
	var Words []string
	re := regexp.MustCompile(`[\x60 1234567890\-=/*\[\];',./~!@#$%^&*()_+{}:"<>?·【】；‘，。/！@#￥%……&*（）——：“《》？]+`)
	runeText := []rune(re.ReplaceAllString(text, ""))

	textLen := len(runeText)

	for i := 0; i < textLen; {
		node := trie.root
		j := i

		for ; j < textLen; j++ {
			char := unicode.ToLower(runeText[j]) // 将匹配字符转换为小写
			if child, ok := node.children[char]; ok {
				node = child
				if node.isEnd {
					Words = append(Words, string(runeText[i:j+1]))
				}
			} else {
				break
			}
		}

		// 如果未匹配到敏感词，继续下一个字符
		if j == i {
			i++
		} else {
			i = j
		}
	}

	return SearchResult{
		UserID:      userID,
		CurrentTime: time.Now(),
		Text:        text,
		Words:       Words,
	}
}

// FindWordsNoUserId 返回文本中包含的敏感词列表（忽略大小写）
func (trie *WordsTrie) FindWordsNoUserId(text string) SearchResult {
	var Words []string
	re := regexp.MustCompile(`[\x60 1234567890\-=/*\[\];',./~!@#$%^&*()_+{}:"<>?·【】；‘，。/！@#￥%……&*（）——：“《》？]+`)
	runeText := []rune(re.ReplaceAllString(text, ""))

	textLen := len(runeText)

	for i := 0; i < textLen; {
		node := trie.root
		j := i

		for ; j < textLen; j++ {
			char := unicode.ToLower(runeText[j]) // 将匹配字符转换为小写
			if child, ok := node.children[char]; ok {
				node = child
				if node.isEnd {
					Words = append(Words, string(runeText[i:j+1]))
				}
			} else {
				break
			}
		}

		// 如果未匹配到敏感词，继续下一个字符
		if j == i {
			i++
		} else {
			i = j
		}
	}

	return SearchResult{
		CurrentTime: time.Now(),
		Text:        text,
		Words:       Words,
	}
}

func InitWordFilter() {
	var err error
	// 从文件中加载敏感词并创建字典树
	WorldFilter, err = NewWordsTrieFromFile(constant.WordFilterFilePath)
	if err != nil {
		zap.S().Info("无法加载敏感词文件:", err)
		return
	}
	zap.S().Infof("敏感词过滤器初始化成功")
}
