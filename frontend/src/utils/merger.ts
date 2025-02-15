import {Diff, DIFF_DELETE, DIFF_EQUAL, DIFF_INSERT, diff_match_patch} from 'diff-match-patch';

const dmp = new diff_match_patch();

export function mergeCode(localCode: string, incomingCode: string): string {
    const diffs = generateDiffs(localCode, incomingCode)
    return applyDiffs(localCode, diffs)
}

export function generateDiffs(localCode: string, incomingCode: string): Diff[] {
    const diffs = dmp.diff_main(localCode, incomingCode);
    dmp.diff_cleanupSemantic(diffs);
    return diffs;
}

export function applyDiffs(localCode: string, diffs: Diff[]): string {
    let result = '';
    let localIndex = 0;

    for (const [operation, text] of diffs) {
        switch (operation) {
            case DIFF_EQUAL:
                result += localCode.slice(localIndex, localIndex + text.length);
                localIndex += text.length;
                break;

            case DIFF_DELETE:
                localIndex += text.length;
                break;

            case DIFF_INSERT:
                result += text;
                break;
        }
    }

    return result;
}
