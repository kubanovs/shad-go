# Code Review: gitfame

## Итог

Задача решена корректно — основная логика работает. Ниже описаны места, которые стоит исправить или улучшить.

---

## 1. `comparator.go` — `return 1` вместо `return 0` при равных элементах

Когда все поля A и B совпадают, функция возвращает `1` вместо `0`. Это нарушает транзитивность — `sort.Slice` может работать непредсказуемо или даже вызвать панику.

```go
// было
func (c *comparator) sort(A analyzer.PersonStat, B analyzer.PersonStat) int {
    for _, el := range c.columnOrder {
        ...
        if res != 0 {
            return res
        }
    }
    return 1 // неправильно: равные элементы должны давать 0
}

// должно быть
    return 0
```

---

## 2. `comparator.go` — неверная ветка для равных имён

Когда `aNameLow == bNameLow`, результат `res = -1` (как будто B > A), хотя на самом деле они равны. Нужна явная проверка на равенство.

```go
// было
case "name":
    aNameLow := strings.ToLower(A.Name)
    bNameLow := strings.ToLower(B.Name)

    if aNameLow < bNameLow {
        res = 1
    } else {
        res = -1 // срабатывает и при aNameLow == bNameLow
    }

// должно быть
case "name":
    aNameLow := strings.ToLower(A.Name)
    bNameLow := strings.ToLower(B.Name)

    if aNameLow < bNameLow {
        res = 1
    } else if aNameLow > bNameLow {
        res = -1
    }
    // иначе res остаётся 0 — равные имена
```

---

## 3. `formatter.go` — переменная `json` затеняет пакет `encoding/json`

```go
// было
func jsonFormat(personStats []analyzer.PersonStat) (string, error) {
    ...
    json, err := json.Marshal(personStatsMarshall) // переменная json затеняет пакет json
    if err != nil {
        return "", err
    }
    return string(json), nil
}

// должно быть
func jsonFormat(personStats []analyzer.PersonStat) (string, error) {
    ...
    data, err := json.Marshal(personStatsMarshall)
    if err != nil {
        return "", err
    }
    return string(data), nil
}
```

---

## 4. `gfcmd.go` — глобальные переменные флагов

Все переменные флагов объявлены глобально, что делает тестирование невозможным без сброса глобального состояния.

```go
// было
var (
    revision    string
    repoPath    string
    format      string
    orderBy     string
    useCommiter bool
    extensions  []string
    exclude     []string
    restrictTo  []string
    languages   []string
)

// должно быть — вынести в структуру и биндить к ней
type config struct {
    revision    string
    repoPath    string
    format      string
    orderBy     string
    useCommitter bool
    extensions  []string
    exclude     []string
    restrictTo  []string
    languages   []string
}

func Main() int {
    cfg := &config{}
    rootCmd := &cobra.Command{
        Use:  "gitfame",
        RunE: func(cmd *cobra.Command, args []string) error {
            return process(cfg)
        },
    }
    rootCmd.Flags().StringVar(&cfg.revision, "revision", "HEAD", "revision of git repository")
    // ...остальные флаги аналогично...
}
```

---

## 5. `configs/language.go` — `panic` в конструкторе скрывает причину ошибки

```go
// было
func New() *LanguageConfig {
    languages, err := loadLanguages()
    if err != nil {
        panic("error while read language_extensions.json") // теряем саму ошибку
    }
    return &LanguageConfig{languages: languages}
}

// должно быть — паниковать с реальной ошибкой (файл всегда встроен через embed)
func New() *LanguageConfig {
    languages, err := loadLanguages()
    if err != nil {
        panic(fmt.Sprintf("failed to load language_extensions.json: %v", err))
    }
    return &LanguageConfig{languages: languages}
}
```

---

## 6. Отсутствует валидация `format` до запуска анализа

Сейчас `format` проверяется только внутри `formatter.Format`, уже после того как `git blame` прогонится по всем файлам. Нужно проверять заранее.

```go
// было
func process(cmd *cobra.Command, args []string) error {
    if _, ok := columns[orderBy]; !ok {
        return fmt.Errorf("not valid orderBy arg: %s", orderBy)
    }
    // ... запускаем анализ, а потом узнаём что format некорректный

// должно быть
var validFormats = map[string]struct{}{
    "tabular":    {},
    "csv":        {},
    "json":       {},
    "json-lines": {},
}

func process(cmd *cobra.Command, args []string) error {
    if _, ok := columns[orderBy]; !ok {
        return fmt.Errorf("not valid orderBy arg: %s", orderBy)
    }
    if _, ok := validFormats[format]; !ok {
        return fmt.Errorf("not valid format: %s", format)
    }
    // ... теперь запускаем анализ
}
```

---

## 7. Опечатки в именах

| Файл | Текущее | Правильное |
|------|---------|------------|
| `git.go:29` | `Commiter` | `Committer` |
| `analyzer.go:12` | `useCommiter` | `useCommitter` |
| `git.go:120` | `parceLsTreeOutput` | `parseLsTreeOutput` |

```go
// было
type GitLogObject struct {
    Author     string
    Commiter   string  // опечатка
    CommitHash string
}

// должно быть
type GitLogObject struct {
    Author     string
    Committer  string
    CommitHash string
}
```

---

## 8. `isAnyMatchGlob` — игнорирование ошибки от `filepath.Match`

```go
// было
func isAnyMatchGlob(path string, globs []string) bool {
    for _, glob := range globs {
        if ok, _ := filepath.Match(glob, path); ok { // ошибка ErrBadPattern игнорируется
            return true
        }
    }
    return false
}

// должно быть
func isAnyMatchGlob(path string, globs []string) (bool, error) {
    for _, glob := range globs {
        ok, err := filepath.Match(glob, path)
        if err != nil {
            return false, fmt.Errorf("invalid glob pattern %q: %w", glob, err)
        }
        if ok {
            return true, nil
        }
    }
    return false, nil
}
```

---

## 9. `analyzer.go` — лишний вызов `updateAuthorStats` для пустого файла

После обработки пустого файла `groupsInfo` пустой, но `updateAuthorStats` всё равно вызывается.

```go
// было
if len(groupsInfo) == 0 {
    lastModification, err := git.FileLastLog(...)
    ...
    updateAuthorStatsEmptyFile(lastModification, lsTreeObj.Path, authorStats, a.useCommiter)
}
updateAuthorStats(groupsInfo, lsTreeObj.Path, authorStats) // бесполезный вызов с пустым слайсом

// должно быть
if len(groupsInfo) == 0 {
    lastModification, err := git.FileLastLog(...)
    ...
    updateAuthorStatsEmptyFile(lastModification, lsTreeObj.Path, authorStats, a.useCommiter)
} else {
    updateAuthorStats(groupsInfo, lsTreeObj.Path, authorStats)
}
```

---

## 10. `parseBlameOutput` — строки `previous` и `boundary` попадают в `parseCommitLine`

В porcelain-формате есть строки `previous <hash> <file>` и `boundary`, которые не совпадают ни с одним известным префиксом и падают в `default` → `parseCommitLine`. Ошибка молча игнорируется, но это неявно.

```go
// было
var ignoredPrefixes = []string{
    "author-mail ", "author-time ", "author-tz ", "committer-mail ", "committer-time ", "committer-tz ",
    "summary ",
}

// должно быть — добавить остальные известные префиксы
var ignoredPrefixes = []string{
    "author-mail ", "author-time ", "author-tz ",
    "committer-mail ", "committer-time ", "committer-tz ",
    "summary ", "previous ", "boundary", "filename ",
}
```

> Примечание: если добавить `"filename "` в `ignoredPrefixes`, нужно убрать явную ветку `case strings.HasPrefix(line, "filename ")` из switch или сделать проверку `isIgnoredLine` после case-ов.

---

## 11. Производительность: параллельная обработка файлов

`git blame` вызывается последовательно для каждого файла. При большом репозитории это медленно. Можно использовать пул горутин:

```go
import (
    "sync"
    "sync/atomic"
)

func (a *Analyzer) Analyze() (AnalyzeResult, error) {
    // ... получаем blobFiles как раньше ...

    var (
        mu          sync.Mutex
        authorStats = make(map[string]PersonStat)
        firstErr    error
        processed   atomic.Int64
    )

    sem := make(chan struct{}, 8) // не более 8 параллельных вызовов git blame
    var wg sync.WaitGroup

    for _, file := range blobFiles {
        file := file
        wg.Add(1)
        sem <- struct{}{}

        go func() {
            defer wg.Done()
            defer func() { <-sem }()

            groupsInfo, err := git.Blame(a.repoPath, file.Path, a.revision, a.useCommitter)

            mu.Lock()
            defer mu.Unlock()

            if err != nil && firstErr == nil {
                firstErr = err
                return
            }
            updateAuthorStats(groupsInfo, file.Path, authorStats)

            n := processed.Add(1)
            fmt.Fprintf(os.Stderr, "\r[%d/%d] processed", n, int64(len(blobFiles)))
        }()
    }

    wg.Wait()
    fmt.Fprintln(os.Stderr) // перевод строки после прогресса

    if firstErr != nil {
        return AnalyzeResult{}, firstErr
    }
    // ... формируем результат ...
}
```

---

## Резюме замечаний

| #  | Серьёзность    | Описание                                               |
|----|----------------|--------------------------------------------------------|
| 1  | Баг            | `return 1` вместо `return 0` в comparator             |
| 2  | Баг            | Неверная ветка при равных именах в comparator          |
| 3  | Баг / стиль    | Переменная `json` затеняет пакет `encoding/json`       |
| 4  | Архитектура    | Глобальные переменные флагов                           |
| 5  | Надёжность     | `panic` без реальной ошибки в `language.New()`         |
| 6  | Надёжность     | Нет валидации `format` до запуска анализа              |
| 7  | Стиль          | Опечатки: `Commiter`, `useCommiter`, `parceLsTree`     |
| 8  | Надёжность     | Игнорируется ошибка `filepath.Match`                   |
| 9  | Читаемость     | Лишний вызов `updateAuthorStats` для пустого файла     |
| 10 | Читаемость     | Строки `previous`/`boundary` неявно попадают в default |
| 11 | Производительность | Последовательный `git blame`, можно распараллелить |
