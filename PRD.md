# PRD: Claude Squad Windows 포팅 프로젝트

## 1. 프로젝트 개요
**목표**: claude-squad를 Windows 환경에서 네이티브로 실행 가능하도록 포팅
**기간**: 1-2일
**우선순위**: 높음

## 2. 현재 상황 분석
### 기존 claude-squad 특징
- Unix/Linux 기반 설계
- tmux 의존성
- bash 스크립트 기반 설치
- GitHub CLI 활용

### Windows 환경 제약사항
- tmux 미지원
- bash 스크립트 실행 제한
- 경로 체계 차이 (/ vs \)
- 패키지 관리자 차이

## 3. 기술 요구사항
### 필수 기능
- [x] 다중 AI 에이전트 관리
- [x] 터미널 세션 관리 (tmux 대체)
- [x] GitHub CLI 통합
- [x] 설치 자동화

### Windows 특화 요구사항
- PowerShell 기반 설치 스크립트
- Windows Terminal 활용
- winget/chocolatey 패키지 관리 지원
- cmd/PowerShell 호환성

## 4. 성공 기준
- Windows에서 `cs` 명령어로 실행 가능
- 다중 에이전트 세션 관리 정상 동작
- GitHub CLI와 완전 연동
- 간단한 설치 과정 (1-click 설치)

## 5. 제외 사항
- WSL 의존성 제거 (순수 Windows 환경)
- 복잡한 GUI 개발 제외
- macOS/Linux 호환성 유지 (기존 기능 유지)

## 6. 리스크 및 대응방안
**리스크**: tmux 기능 완전 대체 어려움
**대응**: Windows Terminal의 탭/분할 기능 활용

**리스크**: 설치 과정 복잡성
**대응**: PowerShell 스크립트 자동화